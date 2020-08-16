# Oat

Many languges allow for a way to archive code in a machine readable format. Java has a .jar, Python has a .pyc, GCC has a .h.gch file. Omm also provides an archiver named Oat. Oat stands for Omm actions tree, and can be used by simply typing `oat` into the cmd. 

Lets take a look at how to use this tool.

Imagine we have a large library written in omm. It would not make much sense to have the library user compile the library every time they want to run a new omm script. Oat can allow the library developer to compile it once, and transfer it to his users. 

To compile an omm script to omm, use the `oat build` command. This will generate a file named `<yourfilename>.oat`. To renmae the file in the cmd, you can use the `-o` option. `oat build -o any_file_name.oat`.

To run this oat file, use the `oat run` command. This will call the main function in that oat file. 
