// Author : Ladislav Nevery 2008

#include <winsock2.h>
#include <winsock.h>
#include <wincrypt.h>
#include <windows.h>

#include <commctrl.h>

#pragma comment(lib,"wsock32.lib")
#pragma comment(lib,"crypt32.lib")
#pragma comment(lib,"comctl32.lib")

char* user  =         "root";
BYTE* pasw  =  (BYTE*)"root"; 

HWND  list; BYTE temp[20], resp[20], *chal; 

DWORD WINAPI Sql( void* command , void* onvalue, void* onfield) {
    HCRYPTPROV  prov; HCRYPTHASH hash; DWORD ret=0,no=20; int i,psw=strlen((char*)pasw); static SOCKET s=0;

    SetWindowText(list,(char*)command);

    char* b = (char*)calloc(1<<24,1), *d; // recv buffer is max row size 16 mb

    if(!s) {

        struct sockaddr_in addr = {AF_INET,htons(3306), 127,0,0,1 };  // Put your server Port and IP here
        s = socket(AF_INET,SOCK_STREAM,IPPROTO_TCP);

        if(connect(s,(struct sockaddr*)&addr,sizeof(addr)) <  0 ) return MessageBox(0,  0,"Connect Failed  ",0); 
        i=recv(s,b,1<<24,0);               if (b[4] < 10 ) return MessageBox(0,b+5,"Need MySql > 4.1",0);

        // Read server auth challenge and calc response by making SHA1 hashes from it and password 
        // Details at: http://forge.mysql.com/wiki/MySQL_Internals_ClientServer_Protocol#Handshake_Initialization_Packet

        chal=(BYTE*)b+strlen(b+5)+10; memcpy(chal+8,chal+27,13); 

        if(! CryptAcquireContext(&prov,0,0,PROV_RSA_FULL,0              ) ) 
             CryptAcquireContext(&prov,0,0,PROV_RSA_FULL,CRYPT_NEWKEYSET); 



        // The following source format is quite wide but only one I found that 
        // keeps order in which hashes are applied visible and operations categorized 

        CryptCreateHash(prov,CALG_SHA1,0,0,&hash); CryptHashData(hash,pasw,psw,0); CryptGetHashParam(hash,HP_HASHVAL,temp,&no,0); CryptDestroyHash(hash);
        CryptCreateHash(prov,CALG_SHA1,0,0,&hash); CryptHashData(hash,temp, 20,0); CryptGetHashParam(hash,HP_HASHVAL,resp,&no,0); CryptDestroyHash(hash);
        CryptCreateHash(prov,CALG_SHA1,0,0,&hash); CryptHashData(hash,chal, 20,0);
                                                   CryptHashData(hash,resp, 20,0); CryptGetHashParam(hash,HP_HASHVAL,resp,&no,0); CryptDestroyHash(hash);   
        CryptReleaseContext( prov,0);               
        printf("%s\n", chal);

        // Construct client auth response 
        // Details at: http://forge.mysql.com/wiki/MySQL_Internals_ClientServer_Protocol#Client_Authentication_Packet

                 d = b+4;
          *(int*)d = 1<<2|1<<9|1<<15|1; d+=4;      // features: CLIENT_FOUND_ROWS|CLIENT_PROTOCOL_41|CLIENT_SECURE_CONNECTION|CLIENT_LONG_PASSWORD;
          *(int*)d = 1<<24;             d+=4;      // max packet size = 16Mb
               * d = 8;                 d+=24;     // utf8 charset
          strcpy(d,user);               d+=1 + strlen(user);
               * d = 20;                d+=1;  for(i=0;i<20;i++)  
              d[i] = resp[i]^temp[i];   d+=22;     // XOR encrypt response
          *(int*)b = d-b-4 | 1<<24;                // calc final packet size and id

          send(s,b,  d-b,0);

          recv(s,(char*)&no,4,0); no&=(1<<24)-1;   // in case of login failure server sends us an error text 
        i=recv(s,b,no,0);        if(i==-1||*b)     return MessageBox(0,i==-1?"Timeout":b+3,"Login Failed",0);
    }
    // Send sql command 
    // Details at: http://forge.mysql.com/wiki/MySQL_Internals_ClientServer_Protocol#Command_Packet
    
    d[4]=0x3; strcpy(d+5,(char*)command); *(int*)d=strlen(d+5)+1; i=send(s,d,4+*(int*)d,0);
    
    // Parse and display record set 
    // Details at: http://forge.mysql.com/wiki/MySQL_Internals_ClientServer_Protocol#Result_Set_Header_Packet

    char *p=b, txt[1000]={0}; BYTE typ[1000]={0}; int fields=0, field=0, value=0, row=0, exit=0, rc=0; 

    while (1) {
               rc = 0;     i=recv(s,(char*)&no,4,0);        no&=0xffffff; // This is a bug. server sometimes don't send those 4 bytes together. will fix it 
        while( rc < no && (i=recv(s,b+rc, no-rc ,0)) > 0  ) rc+=i; 
        
        if(i<1) { closesocket(s); s=0;  break; }                                      // Connection lost

        // 0. For non query sql commands we get just single success or failure response
        if(*       b==0x00&&!exit)                                                  break;   // success
        if(*(BYTE*)b==0xff&&!exit)  { b[*(short*)(b+1)+3]=0; MessageBox(0,b+3,0,0); break; } // failure: show server error text

        // 1. first think we receive is number of fields
        if(!fields ) { memcpy(&fields,b,no); field=fields; continue; } 

        // 3. 5. after receiving last field info or row we get this EOF marker
        if (*(BYTE*)b==0xfe && no < 9) if(exit++) break; else continue;        // EOF

        // 4. after receiving all field infos we receive row field values. One row per Receive/Packet
        while( value  ) {
            *txt=0; i=fields-value; __int64 len=1; BYTE g=*(BYTE*)p; 
            g=(g==0||g==251)?0:(g==252)?2:(g==253)?3:(g==254)?8:1; if(g>1)p++; // Specialy packed length 1-8 bytes
            memcpy(&len,p,g); p+=g;

            // Here you can Add support for displaying more DB types like blobs etc.
            if((typ[i]==0xfe||typ[i]==0xfd||typ[i]==0x03||typ[i]==0x08||typ[i]==0xC)) {  // FIELD_TYPE_STRING || FIELD_TYPE_VAR_STRING || FIELD_TYPE_LONG
                if(g) memcpy(txt,p,len); txt[len]=0; 
                typedef long (*TOnValue)(char*,int,int,int);  if(onvalue) ret=((TOnValue)onvalue)(txt,row,i,typ[i]);
            }                       
            p+=len;
            if(!--value) { row++; value=fields; p=b; break; }
        }

        // 2. Second info we get are field infos like name type etc. One field per Receive/Packet
        if( field  ) {
                  i        = fields - field;
            char* cat      =          p; p+=1+*p; *   cat ++=0;
            char* db       =          p; p+=1+*p; *    db ++=0;
            char* table    =          p; p+=1+*p; * table ++=0;
            char* table_   =          p; p+=1+*p; * table_++=0;
            char* name     =          p; p+=1+*p; *  name ++=0; 
            char* name_    =          p; p+=1+*p; *  name_++=0; *p++=0;
            short charset  = *(short*)p; p+=2;
            long  length   = * (long*)p; p+=4;
                  typ[i]   = * (BYTE*)p; p+=1;
            short flags    = *(short*)p; p+=2;
            BYTE  digits   = * (BYTE*)p; p+=3;
            char* Default  =          p; 

            if(!--field) value = fields; p=b; length=max(length*3,60); length=min(length,200);
            typedef long (*TOnField)(char*,int,int,int);  if(onfield) ((TOnField)onfield)(name,row,i,length);
        }       
    } 
    return ret;
}

void OnValue( char* txt, int row, int col, int typ ) {
    LVITEM v={LVIF_TEXT,row,0,0,0,txt};
    if(!col) ListView_InsertItem (list,&v);
    else     ListView_SetItemText(list,row,col,txt);
}

void OnField( char* txt, int row, int col ,int len ) {
    LVCOLUMN c={LVCF_WIDTH|LVCF_TEXT,0,len,txt,col};
             ListView_InsertColumn(list,col,&c);
}

long GetLong( char* txt ) {
    return    atol( txt );
}

int WINAPI WinMain(HINSTANCE  inst, HINSTANCE  prev, LPSTR cmnd, int show) {    
    MSG m; WSADATA wsa; WSAStartup(MAKEWORD(1,1),&wsa); InitCommonControls();

    list=CreateWindow(WC_LISTVIEW,0,WS_VISIBLE|WS_SIZEBOX|WS_MAXIMIZEBOX|WS_SYSMENU|LVS_REPORT,CW_USEDEFAULT,CW_USEDEFAULT,CW_USEDEFAULT,CW_USEDEFAULT,0,0,0,0);   
    ListView_SetExtendedListViewStyle(list ,LVS_EX_FLATSB|LVS_EX_FULLROWSELECT|LVS_EX_HEADERDRAGDROP|LVS_EX_INFOTIP|LVS_EX_ONECLICKACTIVATE|0x10000);
    SetWindowLong( ListView_GetHeader(list),GWL_STYLE,GetWindowLong(ListView_GetHeader(list), GWL_STYLE)^HDS_BUTTONS);   

    long rows = Sql("select count(*) from mysql.user",GetLong, 0);          // example to get no of rows
                Sql("select       *  from mysql.user",OnValue,OnField);  // example to process data

    while(GetMessage(&m,0,0,0)) DispatchMessage( &m); 
    return 0;
}