#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <fcntl.h>
#include <signal.h>

//如果所有管道读端对应的文件描述符被关闭，则write操作会产生信号SIGPIPE
void sighandler(int signo);
int main(void)
{
    int fds[2];
    if(signal(SIGPIPE,sighandler) == SIG_ERR)
    {
        perror("signal error");
        exit(EXIT_FAILURE);
    }
    if(pipe(fds) == -1){
        perror("pipe error");
        exit(EXIT_FAILURE);
    }
    pid_t pid;
    pid = fork();
    if(pid == -1){
        perror("fork error");
        exit(EXIT_FAILURE);
    }
    if(pid == 0){
        close(fds[0]);//子进程关闭读端
        exit(EXIT_SUCCESS);
    }

    close(fds[0]);//父进程关闭读端
    sleep(1);//确保子进程也将读端关闭
    int ret;
    ret = write(fds[1],"hello",5);
    if(ret == -1){
        printf("write error\n");
    }
    return 0;
}

void sighandler(int signo)
{
    printf("catch a SIGPIPE signal and signum = %d\n",signo);
}