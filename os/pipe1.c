/**************
 * readtest.c *
 **************/
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/types.h>
#include <errno.h>
int main()
{
    int pipe_fd[2];
    pid_t pid;
    char r_buf[100];
    char w_buf[4];
    char* p_wbuf;
    int r_num;
    int cmd;

    memset(r_buf,0,sizeof(r_buf));
    memset(w_buf,0,sizeof(r_buf));
    p_wbuf=w_buf;
    if(pipe(pipe_fd)<0) {
        printf("pipe create error\n");
        return -1;
    }

    /*管道两端可分别用描述字fd[0]以及fd[1]来描述，需要注意的是，
    * 管道的两端是固定了任务的。即一端只能用于读，由描述字fd[0]表示，
    * 称其为管道读端；另一端则只能用于写，由描述字fd[1]来表示，称其为管道写端
    */
    if((pid=fork())==0) {
        printf("\n");
        close(pipe_fd[1]); //子进程关闭pipe写端
        sleep(3);
        r_num=read(pipe_fd[0],r_buf,100); //父进程从buffer读取数据
        printf( "read num is %d, the data read from the pipe is %d\n",r_num,atoi(r_buf));

        close(pipe_fd[0]); //读完之后关闭读端
        exit(0);
    } else if(pid>0) {
        close(pipe_fd[0]);//父进程关闭pipe读端
        strcpy(w_buf,"111"); 
        if(write(pipe_fd[1],w_buf,4)!=-1)
            printf("parent write over\n");
        close(pipe_fd[1]);
        printf("parent close fd[1] over\n");
        sleep(10);
    }   
}