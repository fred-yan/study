#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <errno.h>
#include <fcntl.h>

// 阻塞模式时且n>PIPE_BUF：不具有原子性，可能中间有其他进程穿插写入，
// 直到将n字节全部写入才返回，否则阻塞等待写入
#define ERR_EXIT(m) \
        do \
        { \
                perror(m); \
                exit(EXIT_FAILURE); \
        } while(0)

#define TEST_SIZE 68*1024

int main(void)
{
    char a[TEST_SIZE];
    char b[TEST_SIZE];
    char c[TEST_SIZE];

    memset(a, 'A', sizeof(a));
    memset(b, 'B', sizeof(b));
    memset(c, 'C', sizeof(c));

    int pipefd[2];

    int ret = pipe(pipefd);
    if (ret == -1)
        ERR_EXIT("pipe error");

    pid_t pid;
    pid = fork();
    if (pid == 0)//第一个子进程
    {
        close(pipefd[0]);
        ret = write(pipefd[1], a, sizeof(a));
        printf("apid=%d write %d bytes to pipe\n", getpid(), ret);
        exit(0);
    }

    pid = fork();

    
    if (pid == 0)//第二个子进程
    {
        close(pipefd[0]);
        ret = write(pipefd[1], b, sizeof(b));
        printf("bpid=%d write %d bytes to pipe\n", getpid(), ret);
        exit(0);
    }

    pid = fork();

    
    if (pid == 0)//第三个子进程
    {
        close(pipefd[0]);
        ret = write(pipefd[1], c, sizeof(c));
        printf("bpid=%d write %d bytes to pipe\n", getpid(), ret);
        exit(0);
    }


    close(pipefd[1]);
    
    sleep(1);
    int fd = open("test.txt", O_WRONLY | O_CREAT | O_TRUNC, 0644);
    char buf[1024*4] = {0};
    int n = 1;
    while (1)
    {
        ret = read(pipefd[0], buf, sizeof(buf));
        if (ret == 0)
            break;
        printf("n=%02d pid=%d read %d bytes from pipe buf[4095]=%c\n", n++, getpid(), ret, buf[4095]);
        write(fd, buf, ret);

    }
    return 0;    
}