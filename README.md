# 1. What is Lighting RPC anyway?

## 2. What this project intend to address?

concept：

- does rpc needs reflect ? why is that?
- how proto works for rpc ? can we get rid of that?

coding detail：

- how to extract a whole packet from tcp connection?
- Why byte order matters ? `BigEndian`
    - buffer.Write need to know about byte order ? how about binary.Write？
- 

## 3. What happens when rpc call fired?

## 4. Why we need reflect anyway?

所以到底为什么我们在写RPC 框架的时候都在用反射呢？

详细聊聊这两点以及他们的解决方案

### 4.1 调用方不知道服务器提供了哪些函数以及这些函数的具体参数的类型

   首先， 客户端在调用服务端提供的rpc的时候不知道方法有哪些以及他们的参数， 我来举一个例子

你可能已经很熟悉http 格式以及传输方式了， http 里有一个content-type  经常使用 application/json， 这意味着 http 的body 是按照json数据格式传输的（当然你还可以选择其他 xml等），json 的内容是在运行时才能知道的， 也就是说 服务端收到json 的数据之后要 Parse 一下才能知道 里面的内容。 这就是反序列化的过程。

然而json 因为基于文本， 发送之前要转换成字符串等原因 效率比较低， rpc 一般选择pb 或者私有协议 来做。

pb 很好理解， pb  有工具生成桩代码， 这样客户端和服务端都能知道 提供了什么方法， 都有哪些参数， 这些参数的类型都是什么。 而且 二进制发送数据 很高校

私有协议， on the other hand, just like what we did on this project. You have to implement the marshal and unMarshal method yourself.

### 4.2 服务端的序列化反序列化的过程会很难处理 如果没有reflect

想象一下 如果没有 reflect 服务端要处理一个很大很复杂的对象， 这个时候要怎么写入这么多复杂的key ， 对象呢？ 很麻烦， 而且业务逻辑改动 或者增加新的功能的时候 需要写大量的代码来维护 序列化和反序列化的过程， 这是非常痛苦的。

# 5. serialization

## 5.1 Concept

- lighting protocol is just like json or protocol buffer or jce or gob they are same level
- It’s all about serialization and deserialization

### 5.2 Why ByteOrder matters?  —binary.BigEndian

well , you actually don’t need to worry about byte order if u only got one byte. cuz byte order works between the bytes, I mean like

**Example**:

- Consider the integer `0x12345678`, which is 4 bytes:
    - In **big-endian** format, it would be stored in memory as:
        
        ```
        Address:   0x00  0x01  0x02  0x03
        Value:     0x12  0x34  0x56  0x78
        ```
        
    - In **little-endian** format, it would be stored as:
        
        ```
        Address:   0x00  0x01  0x02  0x03
        Value:     0x78  0x56  0x34  0x12
        ```
        
- For a single byte, say `0xAB`, it is simply stored as:
    
    ```
    Address:   0x00
    Value:     0xAB
    ```
    
- There is no order to consider because there is only one byte.