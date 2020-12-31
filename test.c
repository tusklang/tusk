#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <winsock.h>
#include <winsock2.h>

int main(void)
{
    WSADATA wsa;
    WSAStartup(MAKEWORD(2, 2), &wsa);

    fd_set rfds;
    struct timeval tv;
    int retval;

    /* Watch stdin (fd 0) to see when it has input. */

    struct protoent *protoent = getprotobyname("tcp");
    int s = socket(AF_INET, SOCK_STREAM, protoent->p_proto);

    struct sockaddr_in addr_in;
    addr_in.sin_family = 2;
    addr_in.sin_addr.s_addr =
        inet_addr(/* set s_addr based on hostname */
                  inet_ntoa(
                      *(struct in_addr *)*(gethostbyname("example.com")
                                               ->h_addr_list)));
    addr_in.sin_port = htons(80);
    int size = sizeof(addr_in);
    struct sockaddr *addr = (struct sockaddr *)&addr_in;

    connect(s, addr, size);

    send(s, "GET "
            "/ "
            "HTTP"
            "/1.1"
            "\r\nHost: "
            "example.com"
            "\r\nConnection: close"
            "\r\n\r\n",
         strlen("GET "
                "/ "
                "HTTP"
                "/1.1"
                "\r\nHost: "
                "example.com"
                "\r\nConnection: close"
                "\r\n\r\n"),
         0);

    FD_ZERO(&rfds);
    FD_SET(s, &rfds);

    /* Wait up to five seconds. */

    tv.tv_sec = 5;
    tv.tv_usec = 0;

    retval = select(s + 1, &rfds, NULL, NULL, &tv);
    /* Don't rely on the value of tv now! */

    if (retval == -1)
        perror("select()");
    else if (retval)
        printf("Data is available now.\n");
    /* FD_ISSET(0, &rfds) will be true. */
    else
        printf("No data within five seconds.\n");

    exit(EXIT_SUCCESS);
}