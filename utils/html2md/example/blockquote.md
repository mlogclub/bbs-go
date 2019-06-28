系统内置的验证规则如下：


##  格式验证类

 

<blockquote>
        <h3>
            require</h3>
    </blockquote>

 
 验证某个字段必须，例如：
 

<pre><code>&#39;name&#39;=&gt;&#39;require&#39;
</code></pre>

 

<blockquote>
        <h3>
            number 或者 integer</h3>
    </blockquote>

 
 验证某个字段的值是否为数字（采用<code>filter_var</code>验证），例如：
 

<pre><code>&#39;num&#39;=&gt;&#39;number&#39;
</code></pre>

 

<blockquote>
        <h3>
            float</h3>
    </blockquote>

 
 验证某个字段的值是否为浮点数字（采用<code>filter_var</code>验证），例如：
 

<pre><code>&#39;num&#39;=&gt;&#39;float&#39;
</code></pre>

 

<blockquote>
        <h3>
            boolean</h3>
    </blockquote>

 
 验证某个字段的值是否为布尔值（采用<code>filter_var</code>验证），例如：
 

<pre><code>&#39;num&#39;=&gt;&#39;boolean&#39;
</code></pre>

 

<blockquote>
        <h3>
            email</h3>
    </blockquote>

 
 验证某个字段的值是否为email地址（采用<code>filter_var</code>验证），例如：
 

<pre><code>&#39;email&#39;=&gt;&#39;email&#39;
</code></pre>

 

<blockquote>
        <h3>
            array</h3>
    </blockquote>

 
 验证某个字段的值是否为数组，例如：
 

<pre><code>&#39;info&#39;=&gt;&#39;array&#39;
</code></pre>

 

<blockquote>
        <h3>
            accepted</h3>
    </blockquote>

 
 验证某个字段是否为为 yes, on, 或是 1。这在确认&#34;服务条款&#34;是否同意时很有用，例如：
 

<pre><code>&#39;accept&#39;=&gt;&#39;accepted&#39;
</code></pre>

 

<blockquote>
        <h3>
            date</h3>
    </blockquote>

 
 验证值是否为有效的日期，例如：
 

<pre><code>&#39;date&#39;=&gt;&#39;date&#39;
</code></pre>

 
 会对日期值进行<code>strtotime</code>后进行判断。
 

<blockquote>
        <h3>
            alpha</h3>
    </blockquote>

 
 验证某个字段的值是否为字母，例如：
 

<pre><code>&#39;name&#39;=&gt;&#39;alpha&#39;
</code></pre>

 

<blockquote>
        <h3>
            alphaNum</h3>
    </blockquote>

 
 验证某个字段的值是否为字母和数字，例如：
 

<pre><code>&#39;name&#39;=&gt;&#39;alphaNum&#39;
</code></pre>

 

<blockquote>
        <h3>
            alphaDash</h3>
    </blockquote>

 
 验证某个字段的值是否为字母和数字，下划线<code>_</code>及破折号<code>-</code>，例如：
 

<pre><code>&#39;name&#39;=&gt;&#39;alphaDash&#39;
</code></pre>

 

<blockquote>
        <h3>
            chs</h3>
    </blockquote>

 
 验证某个字段的值只能是汉字，例如：
 

<pre><code>&#39;name&#39;=&gt;&#39;chs&#39;
</code></pre>

 

<blockquote>
        <h3>
            chsAlpha</h3>
    </blockquote>

 
 验证某个字段的值只能是汉字、字母，例如：
 

<pre><code>&#39;name&#39;=&gt;&#39;chsAlpha&#39;
</code></pre>

 

<blockquote>
        <h3>
            chsAlphaNum</h3>
    </blockquote>

 
 验证某个字段的值只能是汉字、字母和数字，例如：
 

<pre><code>&#39;name&#39;=&gt;&#39;chsAlphaNum&#39;
</code></pre>

 

<blockquote>
        <h3>
            chsDash</h3>
    </blockquote>

 
 验证某个字段的值只能是汉字、字母、数字和下划线_及破折号-，例如：
 

<pre><code>&#39;name&#39;=&gt;&#39;chsDash&#39;
</code></pre>

 

<blockquote>
        <h3>
            activeUrl</h3>
    </blockquote>

 
 验证某个字段的值是否为有效的域名或者IP，例如：
 

<pre><code>&#39;host&#39;=&gt;&#39;activeUrl&#39;
</code></pre>

 

<blockquote>
        <h3>
            url</h3>
    </blockquote>

 
 验证某个字段的值是否为有效的URL地址（采用<code>filter_var</code>验证），例如：
 

<pre><code>&#39;url&#39;=&gt;&#39;url&#39;
</code></pre>

 

<blockquote>
        <h3>
            ip</h3>
    </blockquote>

 
 验证某个字段的值是否为有效的IP地址（采用<code>filter_var</code>验证），例如：
 

<pre><code>&#39;ip&#39;=&gt;&#39;ip&#39;
</code></pre>

 
 支持验证ipv4和ipv6格式的IP地址。
 

<blockquote>
        <h3>
            dateFormat:format</h3>
    </blockquote>

 
 验证某个字段的值是否为指定格式的日期，例如：
 

<pre><code>&#39;create_time&#39;=&gt;&#39;dateFormat:y-m-d&#39;
</code></pre>

 

##  长度和区间验证类

 

<blockquote>
        <h3>
            in</h3>
    </blockquote>

 
 验证某个字段的值是否在某个范围，例如：
 

<pre><code>&#39;num&#39;=&gt;&#39;in:1,2,3&#39;
</code></pre>

 

<blockquote>
        <h3>
            notIn</h3>
    </blockquote>

 
 验证某个字段的值不在某个范围，例如：
 

<pre><code>&#39;num&#39;=&gt;&#39;notIn:1,2,3&#39;
</code></pre>

 

<blockquote>
        <h3>
            between</h3>
    </blockquote>

 
 验证某个字段的值是否在某个区间，例如：
 

<pre><code>&#39;num&#39;=&gt;&#39;between:1,10&#39;
</code></pre>

 

<blockquote>
        <h3>
            notBetween</h3>
    </blockquote>

 
 验证某个字段的值不在某个范围，例如：
 

<pre><code>&#39;num&#39;=&gt;&#39;notBetween:1,10&#39;
</code></pre>

 

<blockquote>
        <h3>
            length:num1,num2</h3>
    </blockquote>

 
 验证某个字段的值的长度是否在某个范围，例如：
 

<pre><code>&#39;name&#39;=&gt;&#39;length:4,25&#39;
</code></pre>

 
 或者指定长度
 

<pre><code>&#39;name&#39;=&gt;&#39;length:4&#39;
</code></pre>

 

<blockquote>
        <p>
            如果验证的数据是数组，则判断数组的长度。<br style="box-sizing: inherit; -webkit-tap-highlight-color: transparent; text-size-adjust: none; -webkit-font-smoothing: antialiased; display: block; margin: 0.2em;"/>
            如果验证的数据是File对象，则判断文件的大小。</p>
    </blockquote>

 

<blockquote>
        <h3>
            max:number</h3>
    </blockquote>

 
 验证某个字段的值的最大长度，例如：
 

<pre><code>&#39;name&#39;=&gt;&#39;max:25&#39;
</code></pre>

 

<blockquote>
        <p>
            如果验证的数据是数组，则判断数组的长度。<br style="box-sizing: inherit; -webkit-tap-highlight-color: transparent; text-size-adjust: none; -webkit-font-smoothing: antialiased; display: block; margin: 0.2em;"/>
            如果验证的数据是File对象，则判断文件的大小。</p>
    </blockquote>

 

<blockquote>
        <h3>
            min:number</h3>
    </blockquote>

 
 验证某个字段的值的最小长度，例如：
 

<pre><code>&#39;name&#39;=&gt;&#39;min:5&#39;
</code></pre>

 

<blockquote>
        <p>
            如果验证的数据是数组，则判断数组的长度。<br style="box-sizing: inherit; -webkit-tap-highlight-color: transparent; text-size-adjust: none; -webkit-font-smoothing: antialiased; display: block; margin: 0.2em;"/>
            如果验证的数据是File对象，则判断文件的大小。</p>
    </blockquote>

 

<blockquote>
        <h3>
            after:日期</h3>
    </blockquote>

 
 验证某个字段的值是否在某个日期之后，例如：
 

<pre><code>&#39;begin_time&#39; =&gt; &#39;after:2016-3-18&#39;,
</code></pre>

 

<blockquote>
        <h3>
            before:日期</h3>
    </blockquote>

 
 验证某个字段的值是否在某个日期之前，例如：
 

<pre><code>&#39;end_time&#39;   =&gt; &#39;before:2016-10-01&#39;,
</code></pre>

 

<blockquote>
        <h3>
            expire:开始时间,结束时间</h3>
    </blockquote>

 
 验证当前操作（注意不是某个值）是否在某个有效日期之内，例如：
 

<pre><code>&#39;expire_time&#39;   =&gt; &#39;expire:2016-2-1,2016-10-01&#39;,
</code></pre>

 

<blockquote>
        <h3>
            allowIp:allow1,allow2,...</h3>
    </blockquote>

 
 验证当前请求的IP是否在某个范围，例如：
 

<pre><code>&#39;name&#39;   =&gt; &#39;allowIp:114.45.4.55&#39;,
</code></pre>

 
 该规则可以用于某个后台的访问权限
 

<blockquote>
        <h3>
            denyIp:allow1,allow2,...</h3>
    </blockquote>

 
 验证当前请求的IP是否禁止访问，例如：
 

<pre><code>&#39;name&#39;   =&gt; &#39;denyIp:114.45.4.55&#39;,
</code></pre>

 

##  字段比较类

 

<blockquote>
        <h3>
            confirm</h3>
    </blockquote>

 
 验证某个字段是否和另外一个字段的值一致，例如：
 

<pre><code>&#39;repassword&#39;=&gt;&#39;require|confirm:password&#39;
</code></pre>

 
 <code>5.0.4+</code>版本开始，增加了字段自动匹配验证规则，如password和password_confirm是自动相互验证的，只需要使用
 

<pre><code>&#39;password&#39;=&gt;&#39;require|confirm&#39;
</code></pre>

 
 会自动验证和password_confirm进行字段比较是否一致，反之亦然。
 

<blockquote>
        <h3>
            different</h3>
    </blockquote>

 
 验证某个字段是否和另外一个字段的值不一致，例如：
 

<pre><code>&#39;name&#39;=&gt;&#39;require|different:account&#39;
</code></pre>

 

<blockquote>
        <h3>
            eq 或者 = 或者 same</h3>
    </blockquote>

 
 验证是否等于某个值，例如：
 

<pre><code>&#39;score&#39;=&gt;&#39;eq:100&#39;
&#39;num&#39;=&gt;&#39;=:100&#39;
&#39;num&#39;=&gt;&#39;same:100&#39;
</code></pre>

 

<blockquote>
        <h3>
            egt 或者 &gt;=</h3>
    </blockquote>

 
 验证是否大于等于某个值，例如：
 

<pre><code>&#39;score&#39;=&gt;&#39;egt:60&#39;
&#39;num&#39;=&gt;&#39;&gt;=:100&#39;
</code></pre>

 

<blockquote>
        <h3>
            gt 或者 &gt;</h3>
    </blockquote>

 
 验证是否大于某个值，例如：
 

<pre><code>&#39;score&#39;=&gt;&#39;gt:60&#39;
&#39;num&#39;=&gt;&#39;&gt;:100&#39;
</code></pre>

 

<blockquote>
        <h3>
            elt 或者 &lt;=</h3>
    </blockquote>

 
 验证是否小于等于某个值，例如：
 

<pre><code>&#39;score&#39;=&gt;&#39;elt:100&#39;
&#39;num&#39;=&gt;&#39;&lt;=:100&#39;
</code></pre>

 

<blockquote>
        <h3>
            lt 或者 &lt;</h3>
    </blockquote>

 
 验证是否小于某个值，例如：
 

<pre><code>&#39;score&#39;=&gt;&#39;lt:100&#39;
&#39;num&#39;=&gt;&#39;&lt;:100&#39;
</code></pre>

 

<blockquote>
        <h3>
            验证字段比较支持对比其他字段（V5.0.8+）</h3>
    </blockquote>

 
 验证对比其他字段大小（数值大小对比），例如：
 

<pre><code>&#39;price&#39;=&gt;&#39;lt:market_price&#39;
&#39;price&#39;=&gt;&#39;&lt;:market_price&#39;
</code></pre>

 

##  filter验证


 支持使用filter_var进行验证，例如：
 

<pre><code>&#39;ip&#39;=&gt;&#39;filter:validate_ip&#39;
</code></pre>

 

##  正则验证


 支持直接使用正则验证，例如：
 

<pre><code>&#39;zip&#39;=&gt;&#39;\d{6}&#39;,
// 或者
&#39;zip&#39;=&gt;&#39;regex:\d{6}&#39;,
</code></pre>

 
 如果你的正则表达式中包含有<code>|</code>符号的话，必须使用数组方式定义。
 

<pre><code>&#39;accepted&#39;=&gt;[&#39;regex&#39;=&gt;&#39;/^(yes|on|1)$/i&#39;],
</code></pre>

 
 也可以实现预定义正则表达式后直接调用，例如在验证器类中定义<code>regex</code>属性
 

<pre><code>    protected $regex = [ &#39;zip&#39; =&gt; &#39;\d{6}&#39;];
</code></pre>

 
 然后就可以使用
 

<pre><code>&#39;zip&#39;	=&gt;	&#39;regex:zip&#39;,
</code></pre>

 

##  上传验证

 

<blockquote>
        <h3>
            file</h3>
    </blockquote>

 
 验证是否是一个上传文件
 

<blockquote>
        <h3>
            image:width,height,type</h3>
    </blockquote>

 
 验证是否是一个图像文件，width height和type都是可选，width和height必须同时定义。
 

<blockquote>
        <h3>
            fileExt:允许的文件后缀</h3>
    </blockquote>

 
 验证上传文件后缀
 

<blockquote>
        <h3>
            fileMime:允许的文件类型</h3>
    </blockquote>

 
 验证上传文件类型
 

<blockquote>
        <h3>
            fileSize:允许的文件字节大小</h3>
    </blockquote>

 
 验证上传文件大小


##  行为验证


 使用行为验证数据，例如：
 

<pre><code>&#39;data&#39;=&gt;&#39;behavior:\app\index\behavior\Check&#39;
</code></pre>

 

##  其它验证

 

<blockquote>
        <h3>
            unique:table,field,except,pk</h3>
    </blockquote>

 
<table><thead><tr><th> 版本</th><th> 调整功能</th></tr></thead><tbody><tr><td> 5.0.5</td><td> 支持指定完整模型类 并且默认会优先检测模型类是否存在 不存在则检测数据表</td></tr></tbody></table>

 验证当前请求的字段值是否为唯一的，例如：
 

<pre><code>// 表示验证name字段的值是否在user表（不包含前缀）中唯一
&#39;name&#39;   =&gt; &#39;unique:user&#39;,
// 验证其他字段
&#39;name&#39;   =&gt; &#39;unique:user,account&#39;,
// 排除某个主键值
&#39;name&#39;   =&gt; &#39;unique:user,account,10&#39;,
// 指定某个主键值排除
&#39;name&#39;   =&gt; &#39;unique:user,account,10,user_id&#39;,
</code></pre>

 
 如果需要对复杂的条件验证唯一，可以使用下面的方式：
 

<pre><code>// 多个字段验证唯一验证条件
&#39;name&#39;   =&gt; &#39;unique:user,status^account&#39;,
// 复杂验证条件
&#39;name&#39;   =&gt; &#39;unique:user,status=1&amp;account=&#39;.$data[&#39;account&#39;],
</code></pre>

 

<blockquote>
        <h3>
            requireIf:field,value</h3>
    </blockquote>

 
 验证某个字段的值等于某个值的时候必须，例如：
 

<pre><code>// 当account的值等于1的时候 password必须
&#39;password&#39;=&gt;&#39;requireIf:account,1&#39;
</code></pre>

 

<blockquote>
        <h3>
            requireWith:field</h3>
    </blockquote>

 
 验证某个字段有值的时候必须，例如：
 

<pre><code>// 当account有值的时候password字段必须
&#39;password&#39;=&gt;&#39;requireWith:account&#39;</code></pre>

