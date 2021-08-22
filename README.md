# sshc

免登陆的SSH客户端，支持执行命令和上传下载文件。

可以方便与其他UNIX命令协同。

## 安装方法

1. 克隆本仓库：`git clone https://github.com/daqing/sshc.git`
2. 执行命令：`cd sshc && go install`

## 使用方法

经过安装后，我们得到了`sshc`命令。

它的运行原理是，通过读取配置文件，直接登陆服务器运行命令，免去普通 ssh
客户端需要交互式的输入密码而导致无法自动化的问题。

### 配置文件的格式

    ```
    Key="/home/ubuntu/.ssh/id_rsa"
    IP="cloud.example.com"
    User="root"
    Password=""
    ```

这个配置，适合：服务器只能通过key来登录，禁用了密码登录的情况。
所以这里`Password`要留空

    ```
    Key=""
    IP="cloud.example.com"
    User="root"
    Password="foobar2000"
    ```

这个配置，适合：服务器使用普通的用户名和密码验证的方式登录。
所以这里`Key`要留空。

### 调用参数

  `sshc -c [path/to/config/dir] [host] [action] [arg1 arg2 arg3...]`

### 举例说明

  1. 假设你在`$HOME`目录创建一个文件夹`hosts`，并在`hosts`目录创建一个`cloud.toml`的配置文件
  2. `$HOME/hosts/cloud.toml`文件内容如下：
  
      ```toml
      Key=""
      IP="cloud.example.com"
      User="root"
      Password="foobar2000"
      ```
      
      字段解释：
      
        * Key: 对于普通的使用密码登录的情况，Key留空
        * IP: 我们要登录的服务器IP（或域名）
        * User: 用于ssh登录的用户名，比如`root`或`ubuntu`
        * Password: ssh 登录密码
      
  3. 那么，要自动登录服务器并获取一个命令的返回结果（标准输出），可以运行以下命令：
  
      `sshc -c $HOME/hosts cloud run "ls /tmp"`
      
      这样会自动登录`cloud`服务器，执行`ls /tmp`命令，并在本地终端显示远程命令的标准输出
      
  4. 如果你觉得每次输入那么多参数很麻烦，可以使用shell的`alias`功能。
  
      以`zsh`为例：
      
        `alias scloud="sshc -c $HOME/hosts cloud"`
      
      这样，你只要执行：
      
        `scloud run "ls /tmp"`
      
      就可以了。
      
      `bash`同理。


## 改进与建议

如果你对本项目有好的想法，或者改进建议，请直接提交issue。

如果你懂编程，可以`fork`本项目，实现你想要的功能，
然后提交`pull request`给我们就可以了。

