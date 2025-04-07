# Create 
create VNAME TYPE
```
create c [10][10]i32
create b string 
create d i32 
create e []char
create f f32
```
- All the possible types (TYPE)
    - Primitives (PTYPE)
        - i/u8-i/u64
        - f16-f128
        - cf32-cf256
        - char 
        - bool
    - Arrays of primitive types (ATYPE)
        - Arbitrary length
            - []TYPE
        - Strinct length 
            - `[10][10]TYPE`
    - Structures (STYPE)
        - `struct { NAME TYPE, NAME TYPE, ... }
    - Enums (ETYPE)
        - `enum { NAME1, NAME2, NAME3, ... }
    - Unions (UTYPE)
        - `union { NAME TYPE, NAME TYPE, ... }

# Write 
write WV WFMT LEN DATA
```
write a 10
write [b, c, (d, e)] [1, 2, 3]
```

# Read
read VNAME RFMT 

# Take

# Open
open write aa VNAME FMT 
write aa LEN DATA
write aa LEN DATA
write aa LEN DATA

open read bb VNAME RFMT
read bb LEN 
read bb LEN

# Delete 
delete a
delete a\[0:10\]

# Invariants 
- Len(read) % size(elem) == 0
