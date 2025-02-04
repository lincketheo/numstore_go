let SessionLoad = 1
let s:so_save = &g:so | let s:siso_save = &g:siso | setg so=0 siso=0 | setl so=-1 siso=-1
let v:this_session=expand("<sfile>:p")
silent only
silent tabonly
cd ~/Development/numstore/ndbgo/c
if expand('%') == '' && !&modified && line('$') <= 1 && getline(1) == ''
  let s:wipebuf = bufnr('%')
endif
let s:shortmess_save = &shortmess
if &shortmess =~ 'A'
  set shortmess=aoOA
else
  set shortmess=aoO
endif
badd +408 apps/tests.c
badd +26 ~/.config/nvim/init.vim
badd +50 src/utils.c
badd +13 include/utils.h
badd +20 term://~/Development/numstore/ndbgo/c//460800:/usr/bin/zsh
badd +1 src/ndbc_write.c
argglobal
%argdel
$argadd NvimTree_1
edit apps/tests.c
let s:save_splitbelow = &splitbelow
let s:save_splitright = &splitright
set splitbelow splitright
wincmd _ | wincmd |
vsplit
wincmd _ | wincmd |
vsplit
2wincmd h
wincmd w
wincmd w
let &splitbelow = s:save_splitbelow
let &splitright = s:save_splitright
wincmd t
let s:save_winminheight = &winminheight
let s:save_winminwidth = &winminwidth
set winminheight=0
set winheight=1
set winminwidth=0
set winwidth=1
exe 'vert 1resize ' . ((&columns * 30 + 95) / 190)
exe 'vert 2resize ' . ((&columns * 79 + 95) / 190)
exe 'vert 3resize ' . ((&columns * 79 + 95) / 190)
argglobal
enew
file NvimTree_1
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal nofen
lcd ~/Development/numstore/ndbgo/c
wincmd w
argglobal
balt ~/Development/numstore/ndbgo/c/src/ndbc_write.c
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
31,36fold
39,44fold
47,53fold
56,61fold
64,82fold
85,126fold
140,157fold
160,183fold
195,229fold
232,263fold
274,291fold
294,318fold
322,338fold
342,356fold
360,365fold
369,406fold
414,417fold
420,429fold
436,439fold
442,444fold
409,450fold
let &fdl = &fdl
let s:l = 408 - ((137 * winheight(0) + 25) / 50)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 408
normal! 05|
lcd ~/Development/numstore/ndbgo/c
wincmd w
argglobal
if bufexists(fnamemodify("term://~/Development/numstore/ndbgo/c//460800:/usr/bin/zsh", ":p")) | buffer term://~/Development/numstore/ndbgo/c//460800:/usr/bin/zsh | else | edit term://~/Development/numstore/ndbgo/c//460800:/usr/bin/zsh | endif
if &buftype ==# 'terminal'
  silent file term://~/Development/numstore/ndbgo/c//460800:/usr/bin/zsh
endif
balt ~/Development/numstore/ndbgo/c/src/utils.c
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
let s:l = 1 - ((0 * winheight(0) + 25) / 50)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 1
normal! 020|
lcd ~/Development/numstore/ndbgo/c
wincmd w
2wincmd w
exe 'vert 1resize ' . ((&columns * 30 + 95) / 190)
exe 'vert 2resize ' . ((&columns * 79 + 95) / 190)
exe 'vert 3resize ' . ((&columns * 79 + 95) / 190)
tabnext 1
if exists('s:wipebuf') && len(win_findbuf(s:wipebuf)) == 0 && getbufvar(s:wipebuf, '&buftype') isnot# 'terminal'
  silent exe 'bwipe ' . s:wipebuf
endif
unlet! s:wipebuf
set winheight=1 winwidth=20
let &shortmess = s:shortmess_save
let &winminheight = s:save_winminheight
let &winminwidth = s:save_winminwidth
let s:sx = expand("<sfile>:p:r")."x.vim"
if filereadable(s:sx)
  exe "source " . fnameescape(s:sx)
endif
let &g:so = s:so_save | let &g:siso = s:siso_save
set hlsearch
doautoall SessionLoadPost
unlet SessionLoad
" vim: set ft=vim :
