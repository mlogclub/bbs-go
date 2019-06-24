

# [1. 程序构建](#id11)[¶](#program-build)



目录

- [程序构建](#program-build)
  - [配置](#id3)
  - [编译](#id4)
    - [makefile编写的要点](#makefile)
    - [makefile中的全局自变量](#id5)
    - [更多选择 CMake](#cmake)
    - [编译依赖的库](#id7)
    - [g++编译](#g)
  - [安装](#id8)
  - [总结](#id9)

一般源代码提供的程序安装需要通过配置、编译、安装三个步骤；

- 配置做的工作主要是检查当前环境是否满足要安装软件的依赖关系，以及设置程序安装所需要的初始化信息，比如安装路径，需要安装哪些组件；配置完成，会生成makefile文件供第二步make使用；
- 编译是对源文件进行编译链接生成可执行程序；
- 安装做的工作就简单多了，就是将生成的可执行文件拷贝到配置时设置的初始路径下；


## [1.1. 配置](#id12)[¶](#id3)


查询可用的配置选项:




<pre>./configure --help
</pre>



配置路径:




<pre>./configure --prefix=/usr/local/snmp
</pre>



–prefix是配置使用的最常用选项，设置程序安装的路径；




## [1.2. 编译](#id13)[¶](#id4)


编译使用make编译:




<pre>make -f myMakefile
</pre>



通过-f选项显示指定需要编译的makefile；如果待使用makefile文件在当前路径，且文件名为以下几个，则不用显示指定：

makefile Makefile



### [makefile编写的要点](#id14)[¶](#makefile)


- 必须满足第一条规则，满足后停止
- 除第一条规则，其他无顺序



### [makefile中的全局自变量](#id15)[¶](#id5)


- $@目标文件名
- @^所有前提名，除副本
- @＋所有前提名，含副本
- @＜一个前提名
- @？所有新于目标文件的前提名
- @*目标文件的基名称

注解

系统学习makefile的书写规则，请参考 跟我一起学makefile [[1]](#id10)





### [更多选择 CMake](#id16)[¶](#cmake)


CMake是一个跨平台的安装（编译）工具，可以用简单的语句来描述所有平台的安装(编译过程)。他能够输出各种各样的makefile或者project文件。使用CMake，能够使程序员从复杂的编译连接过程中解脱出来。它使用一个名为 CMakeLists.txt 的文件来描述构建过程,可以生成标准的构建文件,如 Unix/Linux 的 Makefile 或Windows Visual C++ 的 projects/workspaces 。




### [编译依赖的库](#id17)[¶](#id7)


makefile编译过程中所依赖的非标准库和头文件路径需要显示指明:




<pre>CPPFLAGS -I标记非标准头文件存放路径
LDFLAGS  -L标记非标准库存放路径
</pre>



如果CPPFLAGS和LDFLAGS已在用户环境变量中设置并且导出（使用export关键字），就不用再显示指定；




<pre>make -f myMakefile LDFLAGS=&#39;-L/var/xxx/lib -L/opt/mysql/lib&#39;
    CPPFLAGS=&#39;-I/usr/local/libcom/include -I/usr/local/libpng/include&#39;
</pre>




警告

链接多库时，多个库之间如果有依赖，需要注意书写的顺序，右边是左边的前提；





### [g++编译](#id18)[¶](#g)





<pre>g++ -o unixApp unixApp.o a.o b.o
</pre>



选项说明：

- -o:指明生成的目标文件
- -g：添加调试信息
- -E: 查看中间文件
应用：查询宏展开的中间文件：

在g++的编译选项中，添加 -E选项，然后去掉-o选项 ，重定向到一个文件中即可:




<pre>g++ -g -E unixApp.cpp  -I/opt/app/source &gt; midfile
</pre>



查询应用程序需要链接的库:




<pre>$ldd myprogrammer
    libstdc++.so.6 =&gt; /usr/lib64/libstdc++.so.6 (0x00000039a7e00000)
    libm.so.6 =&gt; /lib64/libm.so.6 (0x0000003996400000)
    libgcc_s.so.1 =&gt; /lib64/libgcc_s.so.1 (0x00000039a5600000)
    libc.so.6 =&gt; /lib64/libc.so.6 (0x0000003995800000)
    /lib64/ld-linux-x86-64.so.2 (0x0000003995400000)
</pre>




注解

关于ldd的使用细节，参见 [ldd 查看程序依赖库](../tool/ldd.html#ldd)






## [1.3. 安装](#id19)[¶](#id8)


安装做的工作就简单多了，就是将生成的可执行文件拷贝到配置时设置的初始路径下:




<pre>$make install
</pre>



其实 **install** 就是makefile中的一个规则，打开makefile文件后可以查看程序安装的所做的工作；




## [1.4. 总结](#id20)[¶](#id9)


configure make install g++


<table rules="none"><colgroup><col class="label"/><col/></colgroup><tbody valign="top"><tr><td>[[1]](#id6)</td><td>陈皓 跟我一起写Makefile [http://scc.qibebt.cas.cn/docs/linux/base/%B8%FA%CE%D2%D2%BB%C6%F0%D0%B4Makefile-%B3%C2%F0%A9.pdf](http://scc.qibebt.cas.cn/docs/linux/base/%B8%FA%CE%D2%D2%BB%C6%F0%D0%B4Makefile-%B3%C2%F0%A9.pdf)</td></tr></tbody></table>




