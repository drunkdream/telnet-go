# Telnet of go

这是一个用go实现的`telnet`程序，你可以把它当作一个普通的telnet客户端来用（访问中文telnet服务端可能会有乱码）。当然，它的真正用途并不在此，而是用于当SSH服务端不支持端口转发时建立一个TCP隧道。实现原理是通过将socket双向通信转换为对`stdin`和`stdout`的读写，而`stderr`则用于日志或错误信息的输出。


