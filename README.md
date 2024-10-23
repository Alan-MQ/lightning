# What is Lighting RPC anyway?

## What this project intend to address?

concept：

- does rpc needs reflect ? why is that?
- how proto works for rpc ? can we get rid of that?

coding detail：

- how to extract a whole packet from tcp connection?
- Why byte order matters ? `BigEndian`
    - buffer.Write need to know about byte order ? how about binary.Write？
- 

## What happens when rpc call fired?

## Detail of serialization

### Why ByteOrder matters?  —binary.BigEndian

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