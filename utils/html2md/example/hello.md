###### 3.5.2.1.43. 表格


[在线示例](https://demo.cuba-platform.com/sampler/open?screen=simple-table)

[API 文档](http://files.cuba-platform.com/javadoc/cuba/7.0/com/haulmont/cuba/gui/components/Table.html)


<code>Table</code> 组件以表格的方式展示信息，对数据进行排序 、管理表格列和表头，并对选中的行执行操作。



![gui table](./img/gui_table.png)



组件对应的 XML 名称： <code>table</code>



在界面 XML 描述中定义组件的示例：



 

<pre><code>&lt;data readOnly=&#34;true&#34;&gt;
    &lt;collection id=&#34;ordersDc&#34; class=&#34;com.company.sales.entity.Order&#34; view=&#34;order-with-customer&#34;&gt;
        &lt;loader id=&#34;ordersDl&#34;&gt;
            &lt;query&gt;
                &lt;![CDATA[select e from sales_Order e]]&gt;
            &lt;/query&gt;
        &lt;/loader&gt;
    &lt;/collection&gt;
&lt;/data&gt;
&lt;layout&gt;
&lt;table id=&#34;ordersTable&#34; dataContainer=&#34;ordersDc&#34; width=&#34;100%&#34;&gt;
    &lt;columns&gt;
        &lt;column id=&#34;date&#34;/&gt;
        &lt;column id=&#34;amount&#34;/&gt;
        &lt;column id=&#34;customer&#34;/&gt;
    &lt;/columns&gt;
    &lt;rowsCount/&gt;
&lt;/table&gt;
&lt;/layout&gt;</code></pre>

 



在上面的示例中，<code>data</code> 元素定义[集合数据容器](gui_collection_container.html)，它使用 [JPQL](glossary.html#jpql) 查询 <code>Order</code> 实体。<code>table</code> 元素定义数据容器，而 <code>columns</code> 元素定义哪些实体属性用作表格列。



<code>table</code> 元素:


[](#gui_Table_rows)
- 
<code>rows</code> – 如果使用 [datasource](gui_attributes.html#gui_attr_datasource) 属性来做数据绑定，则必须设置此元素。


每行可以在左侧的附加列中有一个图标。在界面控制器中创建 <code>ListComponent.IconProvider</code> 接口的实现，并将其设置给表格：



 

<pre><code>@Inject
private Table&lt;Customer&gt; table;

@Subscribe
protected void onInit(InitEvent event) {
    table.setIconProvider(new ListComponent.IconProvider&lt;Customer&gt;() {
        @Nullable
        @Override
        public String getItemIcon(Customer entity) {
            CustomerGrade grade = entity.getGrade();
            switch (grade) {
                case PREMIUM: return &#34;icons/premium_grade.png&#34;;
                case HIGH: return &#34;icons/high_grade.png&#34;;
                case MEDIUM: return &#34;icons/medium_grade.png&#34;;
                default: return null;
            }
        }
    });
}</code></pre>

 



[](#gui_Table_columns)
- 
<code>columns</code> – 定义表格列的必须元素。


每个列都在嵌套的 <code>column</code> 元素中描述，<code>column</code> 元素具有以下属性：




[](#gui_Table_column_id)
  - 
<code>id</code> − 必须属性，包含列中要显示的实体属性的名称。可以是来自数据容器的实体的属性，也可以是关联实体的属性（使用 &#34;.&#34; 来指定属性在对象关系图中的路径）。例如：


 

<pre><code>&lt;columns&gt;
    &lt;column id=&#34;date&#34;/&gt;
    &lt;column id=&#34;customer&#34;/&gt;
    &lt;column id=&#34;customer.name&#34;/&gt;
    &lt;column id=&#34;customer.address.country&#34;/&gt;
&lt;/columns&gt;</code></pre>

 



[](#gui_Table_column_caption)
  - 
<code>caption</code> − 包含列标题的可选属性。如果不指定，将显示[本地化属性名称](entity_localization.html)。


[](#gui_Table_column_collapsed)
  - 
<code>collapsed</code> − 可选属性；当设置为 <code>true</code> 时默认隐藏列。当表格的 <code>columnControlVisible</code> 属性不是 <code>false</code> 时，用户可以通过表格右上角的菜单中的按钮 ![gui_table_columnControl](./img/gui_table_columnControl.png) 控制列的可见性。默认情况下，<code>collapsed</code> 是 <code>false</code>。


[](#gui_Table_column_width)
  - 
<code>width</code> − 可选属性，控制默认列宽。只能是以像素为单位的数值。


[](#gui_Table_column_align)
  - 
<code>align</code> − 可选属性，用于设置单元格的文本对齐方式。可选值：<code>LEFT</code> 、 <code>RIGHT</code> 、 <code>CENTER</code>。默认为 <code>LEFT</code>。


[](#gui_Table_column_editable)
  - 
<code>editable</code> − 可选属性，允许编辑表中的相应列。为了使列可编辑，整个表的 [editable](#gui_Table_editable) 属性也应设置为 <code>true</code>。不支持在运行时更改此属性。


[](#gui_Table_column_sortable)
  - 
<code>sortable</code> − 可选属性，用于禁用列的排序。整个表的 [sortable](#gui_Table_sortable) 属性为 <code>true</code> 此属性有效（默认为 <code>true</code>）。


[](#gui_Table_column_maxTextLength)
  - 
<code>maxTextLength</code> – 可选属性，允许限制单元格中的字符数。如果实际值和最大允许字符数之间的差异不超过 10 个字符，则多出来的字符不会被隐藏。用户可以点击可见部分来查看完整的文本。例如一列的字符数限制为 10 个字符：


![gui table column maxTextLength](./img/gui_table_column_maxTextLength.png)



[](#gui_Table_column_link)
  - 
<code>link</code> - 如果设置为 <code>true</code>，则允许列中显示指向实体编辑器的链接。对于基本类型的列，<code>link</code> 属性也可以设置为 true; 在这种情况下，将打开主实体编辑器。这个方法可用于简化导航：用户能够通过单击一些关键属性快速地打开实体编辑器。


[](#gui_Table_column_linkScreen)
  - 
<code>linkScreen</code> - 设置单击 <code>link</code> 属性为 <code>true</code> 的列中的链接时打开的界面的标识符。


[](#gui_Table_column_linkScreenOpenType)
  - 
<code>linkScreenOpenType</code> - 设置界面打开模式（<code>THIS_TAB</code> 、 <code>NEW_TAB</code> 或者 <code>DIALOG</code>）。


[](#gui_Table_column_linkInvoke)
  - 
<code>linkInvoke</code> - 单击链接时调用控制器方法而不是打开界面。


 

<pre><code>@Inject
private Notifications notifications;

public void linkedMethod(Entity item, String columnId) {
    Customer customer = (Customer) item;
    notifications.create()
            .withCaption(customer.getName())
            .show();
}</code></pre>

 



[](#gui_Table_column_captionProperty)
  - 
<code>captionProperty</code> - 指定一个要显示在列中的实体属性名称，而不是显示 [id](#gui_Table_column_id) 指定的实体属性值。例如，如果有一个包含 <code>name</code> 和 <code>orderNo</code> 属性的实体 <code>Priority</code>，则可以定义以下列：


 

<pre><code>&lt;column id=&#34;priority.orderNo&#34; captionProperty=&#34;priority.name&#34; caption=&#34;msg://priority&#34; /&gt;</code></pre>

 



此时，列中将会显示 <code>Priority</code> 实体的 <code>name</code> 属性，但是列的排序是根据 <code>Priority</code> 实体的 <code>orderNo</code> 属性。



[](#gui_Table_column_generator)
  - 
可选的 <code>generator</code> 属性包含指向界面控制器中方法，该方法可创建一个可视化组件显示在表格单元格中：


 

<pre><code>&lt;columns&gt;
    &lt;column id=&#34;name&#34;/&gt;
    &lt;column id=&#34;imageFile&#34;
            generator=&#34;generateImageFileCell&#34;/&gt;
&lt;/columns&gt;</code></pre>

 



 

<pre><code>public Component generateImageFileCell(Employee entity){
    Image image = uiComponents.create(Image.NAME);
    image.setSource(FileDescriptorResource.class).setFileDescriptor(entity.getImageFile());
    return image;
}</code></pre>

 



它可以用来为 [addGeneratedColumn()](#gui_Table_addGeneratedColumn) 方法提供一个 <code>Table.ColumnGenerator</code> 的实现


  - 
<code>column</code> 元素可能包含一个嵌套的[formatter](gui_formatter.html)元素，它允许以不同于[Datatype](datatype.html)的标准格式显示属性值：


 

<pre><code>&lt;column id=&#34;date&#34;&gt;
    &lt;formatter class=&#34;com.haulmont.cuba.gui.components.formatters.DateFormatter&#34;
               format=&#34;yyyy-MM-dd HH:mm:ss&#34;
               useUserTimezone=&#34;true&#34;/&gt;
&lt;/column&gt;</code></pre>

 






[](#gui_Table_rowsCount)
- 
<code>rowsCount</code> − 可选元素，为表格添加 <code>RowsCount</code> 组件；此组件能够分页加载表格数据。可以使用[数据加载器](gui_data_loaders.html)的 <code>setMaxResults()</code> 方法限制数据容器中的记录数来定义分页的大小。这个方法通常是由链接到表格数据加载器的[过滤器](gui_Filter.html)组件来执行的。如果表格没有通用过滤器，则可以直接从界面控制器调用此方法。


<code>RowsCount</code> 组件还可以显示当前数据容器查询的记录总数，而无需提取这些记录。当用户单击 **?** 图标时，它会调用 <code>com.haulmont.cuba.core.global.DataManager#getCount</code> 方法，执行与当前查询条件相同的数据库查询，不过会使用 <code>COUNT(*)</code> 聚合函数代替查询列。然后显示检索到的数字，代替 **?** 图标。



[](#gui_Table_actions)
- 
<code>actions</code> − 可选元素，用于描述和表格相关的[操作](gui_Action.html)。除了自定义的操作外，该元素还支持以下在 <code>com.haulmont.cuba.gui.actions.list</code> 里定义[标准操作](std_list_actions.html)：<code>create</code> 、 <code>edit</code> 、 <code>remove</code> 、 <code>refresh</code> 、 <code>add</code> 、 <code>exclude</code> 、 <code>excel</code>。


[](#gui_Table_buttonsPanel)
- 
可选元素，在表格上方添加一个 [ButtonsPanel](gui_ButtonsPanel.html) 容器来显示操作按钮。



<code>table</code> 属性:


[](#gui_Table_multiselect)
- 
<code>multiselect</code> 属性可以为表格行设置多选模式。如果 <code>multiselect</code> 为 <code>true</code>，用户可以按住 **Ctrl** 或 **Shift** 键在表格中选择多行。默认情况下关闭多选模式。


[](#gui_Table_sortable)
- 
<code>sortable</code> 属性可以对表中的数据进行排序。默认情况下，它设置为 <code>true</code> 。如果允许排序，单击列标题在列名称右侧将显示图标 ![gui_sortable_down](./img/gui_sortable_down.png) / ![gui_sortable_up](./img/gui_sortable_up.png)。可以使用[sortable](#gui_Table_column_sortable)属性禁用特定列的排序。


根据是否将所有记录放在了一页上来使用不同的方式进行排序。如果所有记录在一页，则在内存中执行排序而不需要数据库查询。如果数据有多页，则通过发送具有相应 <code>ORDER BY</code> 条件的新的查询请求在数据库中执行排序。



一个表格中的列可能包含本地属性或实体链接。例如：



 

<pre><code>&lt;table id=&#34;ordersTable&#34;
       dataContainer=&#34;ordersDc&#34;&gt;
    &lt;columns&gt;
        &lt;column id=&#34;customer.name&#34;/&gt; &lt;!-- the &#39;name&#39; attribute of the &#39;Customer&#39; entity --&gt;
        &lt;column id=&#34;contract&#34;/&gt;      &lt;!-- the &#39;Contract&#39; entity --&gt;
    &lt;/columns&gt;
&lt;/table&gt;</code></pre>

 



在后一种情况下，数据排序是根据关联实体的 <code>@NamePattern</code> 注解中定义的属性进行的。如果实体中没有这个注解，则仅仅在内存中对当前页的数据进行排序。



如果列引用了非持久化实体属性，则数据排序将根据 <code>@MetaProperty</code> 注解的 <code>related()</code> 参数中定义的属性执行。如果未指定相关属性，则仅仅在内存中对当前页的数据进行排序。



如果表格链接到一个嵌套的[属性容器](gui_property_containers.html)，这个属性容器包含相关实体的集合。这个集合属性必须是有序类型（<code>List</code> 或 <code>LinkedHashSet</code>）才能使表格支持排序。如果属性的类型为 <code>Set</code>，则 <code>sortable</code> 属性不起作用，并且用户无法对表格进行排序。



[](#gui_Table_presentations)
- 
<code>presentations</code> 属性控制[展示设置](gui_presentations.html)。默认情况下，该值为 <code>false</code>。如果属性值为 <code>true</code>，则会在表格的右上角添加相应的图标 ![gui_presentation](./img/gui_presentation.png)。


[](#gui_Table_columnControlVisible)
- 
如果 <code>columnControlVisible</code> 属性设置为 <code>false</code>，则用户无法使用位于表头的右侧的下拉菜单按钮 ![gui_table_columnControl](./img/gui_table_columnControl.png) 隐藏列：![gui_table_columnControl](./img/gui_table_columnControl.png) 按钮位于表头的右侧。当前显示的列在菜单中标记为选中状态。



![gui table columnControl all](./img/gui_table_columnControl_all.png)


[](#gui_Table_reorderingAllowed)
- 
如果 <code>reorderingAllowed</code> 属性设置为 <code>false</code>，则用户不能通过用鼠标拖动来更改列顺序。


[](#gui_Table_columnHeaderVisible)
- 
如果 <code>columnHeaderVisible</code> 属性设置为 <code>false</code>，则该表没有列标题。


[](#gui_Table_showSelection)
- 
如果 <code>showSelection</code> 属性设置为 <code>false</code>，则不突出显示当前行。


[](#gui_Table_allowPopupMenu)
- 
<code>contextMenuEnabled</code> 属性启用右键菜单。默认情况下，此属性设置为 <code>true</code>。右键菜单中会列出表格操作（如果有的话）和 **System Information** 菜单项（如果用户具有 <code>cuba.gui.showInfo</code> [权限](permissions.html)），通过 **System Information** 菜单项可查看选中实体的详细信息。


[](#gui_Table_multiLineCells)
- 
将 <code>multiLineCells</code> 设置为 <code>true</code> 可以让包含多行文本的单元格显示多行文本。在这种模式下，浏览器会一次加载表格中当前页的所有行，而不是延迟加载表格的可见部分。这就要求在 Web 客户端中适当的滚动。默认值为“false”。


[](#gui_Table_aggregatable)
- 
<code>aggregatable</code> 属性启用表格行的聚合运算。支持以下操作：




  - 
<code>SUM</code> – 计算总和

  - 
<code>AVG</code> – 计算平均值

  - 
<code>COUNT</code> – 计算总数

  - 
<code>MIN</code> – 找到最小值

  - 
<code>MAX</code> – 找到最大值



聚合列的 <code>aggregation</code> 元素应该设置 <code>type</code> 属性，在这个属性中设置聚合函数。默认情况下，聚合列仅支持数值类型，例如 <code>Integer 、 Double 、 Long</code> 和 <code>BigDecimal</code>。聚合表格值显示在表格顶部的附加行中。这是一个定义聚合表示例：



 

<pre><code>&lt;table id=&#34;itemsTable&#34; aggregatable=&#34;true&#34; dataContainer=&#34;itemsDc&#34;&gt;
    &lt;columns&gt;
        &lt;column id=&#34;product&#34;/&gt;
        &lt;column id=&#34;quantity&#34;/&gt;
        &lt;column id=&#34;amount&#34;&gt;
            &lt;aggregation type=&#34;SUM&#34;/&gt;
        &lt;/column&gt;
    &lt;/columns&gt;
&lt;/table&gt;</code></pre>

 



<code>aggregation</code> 元素还可以包含 <code>strategyClass</code> 属性，指定一个实现 <code>AggregationStrategy</code> 接口的类（参阅下面以编程方式设置聚合策略的示例）。



可以指定不同于 [Datatype](datatype.html) 标准格式的[格式化器](gui_formatter.html)显示聚合值：



 

<pre><code>&lt;column id=&#34;amount&#34;&gt;
    &lt;aggregation type=&#34;SUM&#34;&gt;
        &lt;formatter class=&#34;com.company.sample.MyFormatter&#34;/&gt;
    &lt;/aggregation&gt;
&lt;/column&gt;</code></pre>

 



<code>aggregationStyle</code> 属性允许指定聚合行的位置：<code>TOP</code> 或 <code>BOTTOM</code>。默认情况下使用 <code>TOP</code>。



除了上面列出的操作之外，还可以自定义聚合策略，通过实现 <code>AggregationStrategy</code> 接口并将其传递给 <code>AggregationInfo</code> 实例中 <code>Table.Column</code> 类的 <code>setAggregation()</code> 方法。例如：



 

<pre><code>public class TimeEntryAggregation implements AggregationStrategy&lt;List&lt;TimeEntry&gt;, String&gt; {
    @Override
    public String aggregate(Collection&lt;List&lt;TimeEntry&gt;&gt; propertyValues) {
        HoursAndMinutes total = new HoursAndMinutes();
        for (List&lt;TimeEntry&gt; list : propertyValues) {
            for (TimeEntry timeEntry : list) {
                total.add(HoursAndMinutes.fromTimeEntry(timeEntry));
            }
        }
        return StringFormatHelper.getTotalDayAggregationString(total);
    }
    @Override
    public Class&lt;String&gt; getResultClass() {
        return String.class;
    }
}</code></pre>

 



 

<pre><code>AggregationInfo info = new AggregationInfo();
info.setPropertyPath(metaPropertyPath);
info.setStrategy(new TimeEntryAggregation());

Table.Column column = weeklyReportsTable.getColumn(columnId);
column.setAggregation(info);</code></pre>

 





[](#gui_Table_editable)
- 
<code>editable</code> 属性可以将表格转换为即时编辑模式。在这种模式下，具有 <code>editable = true</code> 属性的列显示用于编辑相应实体属性的组件。


根据相应实体属性的类型自动选择每个可编辑列的组件类型。例如，对于字符串和数字属性，应用程序将使用 [TextField](gui_TextField.html)；对于 <code>Date</code> 将使用 [DateField](gui_DateField.html)；对于列表将使用 [LookupField](gui_LookupField.html)；对于指向其它实体的链接将使用 [PickerField](gui_PickerField.html)。



对于 <code>Date</code> 类型的可编辑列，还可以定义 <code>dateFormat</code> 或 <code>resolution</code> 属性，类似于为 [DateField](gui_DateField.html) 的属性。



可以为显示链接实体的可编辑列定义 [optionsContainer](gui_attributes.html#gui_attr_optionsContainer) 和 [captionProperty](gui_attributes.html#gui_attr_captionProperty) 属性。如果设置了 <code>optionsContainer</code> 属性，应用程序将使用 [LookupField](gui_LookupField.html) 而不是 [PickerField](gui_PickerField.html)。



可以使用 <code>Table.addGeneratedColumn()</code> 方法实现单元格的自定义配置（包括编辑） - 见下文。



[](#gui_Table_stylename)
- 
在具有基于 Halo-based 主题的 Web 客户端中，<code>stylename</code> 属性可以在 XML 描述中或者界面控制器中为 <code>Table</code> 组件设置预定义样式：




 

<pre><code>&lt;table id=&#34;table&#34;
       dataContainer=&#34;itemsDc&#34;
       stylename=&#34;no-stripes&#34;&gt;
    &lt;columns&gt;
        &lt;column id=&#34;product&#34;/&gt;
        &lt;column id=&#34;quantity&#34;/&gt;
    &lt;/columns&gt;
&lt;/table&gt;</code></pre>

 



当以编程方式设置样式时，需要选择 <code>HaloTheme</code> 类的一个以 <code>TABLE_</code> 为前缀的常量：



 

<pre><code>table.setStyleName(HaloTheme.TABLE_NO_STRIPES);</code></pre>

 



表格样式:


[](#gui_table_borderless)
  - 
<code>borderless</code> - 不显示表格的外部边线。


[](#gui_table_compact)
  - 
<code>compact</code> - 减少表格单元格内的空白区域。


[](#gui_table_no-header)
  - 
<code>no-header</code> - 隐藏表格的列标题。


[](#gui_table_no-horizontal-lines)
  - 
<code>no-horizontal-lines</code> - 删除行之间的水平分隔线。


[](#gui_table_no-stripes)
  - 
<code>no-stripes</code> - 删除交替的行颜色。


[](#gui_table_no-vertical-lines)
  - 
<code>no-vertical-lines</code> - 删除列之间的垂直分隔线。


[](#gui_table_small)
  - 
<code>small</code> - 使用小字体并减少表格单元格内的空白区域。






<code>Table</code> 接口的方法:


[](#gui_Table_ColumnCollapseListener)
- 
可以使用 <code>addColumnCollapsedListener</code> 方法和 <code>ColumnCollapsedListener</code> 接口的实现跟踪列的可视化状态。


[](#gui_Table_getSelected)
- 
<code>getSelected()</code> 、 <code>getSingleSelected()</code> 返回表格中的选定行对应的实体实例。可以通过调用 <code>getSelected()</code> 方法来获得集合。如果未选择任何内容，则程序将返回空集。如果禁用了 <code>multiselect</code>，应该使用 <code>getSingleSelected()</code> 方法返回一个选定实体，如果没有选择任何内容则返回 <code>null</code>。


[](#gui_Table_SelectionListener)
- 
<code>addSelectionListener()</code> 可以跟踪表格选中行的变化，示例：


 

<pre><code>customersTable.addSelectionListener(customerSelectionEvent -&gt;
        notifications.create()
                .withCaption(&#34;You selected &#34; + customerSelectionEvent.getSelected().size() + &#34; customers&#34;)
                .show());</code></pre>

 



也可以通过订阅相应的事件来跟踪选中行的变化：



 

<pre><code>@Subscribe(&#34;customersTable&#34;)
protected void onCustomersTableSelection(Table.SelectionEvent&lt;Customer&gt; event) {
    notifications.create()
            .withCaption(&#34;You selected &#34; + customerSelectionEvent.getSelected().size() + &#34; customers&#34;)
            .show();
}</code></pre>

 



可以使用[isUserOriginated()](gui_api.html#gui_api_UserOriginated) 方法跟踪 <code>SelectionEvent</code> 事件的来源。



[](#gui_Table_addGeneratedColumn)
- 
<code>addGeneratedColumn()</code> 方法允许在列中自定义数据的表现方式。它需要两个参数：列的标识符和 <code>Table.ColumnGenerator</code> 接口的实现。如果标识符可以匹配 XML 描述中为表格列设置的标识符 - 在这种情况下，插入新列代替 XML 中定义的列。如果标识符与任何列都不匹配，则会在右侧添加新列。




对于表的每一行将调用 <code>Table.ColumnGenerator</code> 接口的 <code>generateCell()</code> 方法。该方法接受在相应行中显示的实体实例作为参数。<code>generateCell()</code> 方法应该返回一个可视化组件，该组件将显示在单元格中。



使用组件的示例：



 

<pre><code>@Inject
private GroupTable&lt;Car&gt; carsTable;
@Inject
private CollectionContainer&lt;Car&gt; carsDc;
@Inject
private CollectionContainer&lt;Color&gt; colorsDc;
@Inject
private UiComponents uiComponents;
@Inject
private Actions actions;

@Subscribe
protected void onInit(InitEvent event) {
    carsTable.addGeneratedColumn(&#34;color&#34;, entity -&gt; {
        LookupPickerField&lt;Color&gt; field = uiComponents.create(LookupPickerField.NAME);
        field.setValueSource(new ContainerValueSource&lt;&gt;(carsDc, &#34;color&#34;));
        field.setOptions(new ContainerOptions&lt;&gt;(colorsDc));
        field.addAction(actions.create(LookupAction.class));
        field.addAction(actions.create(OpenAction.class));
        return field;
    });
}</code></pre>

 



在上面的示例中，表中 <code>color</code> 列中的所有单元格都显示了 [LookupPickerField](gui_LookupPickerField.html) 组件。组件应将它的值保存到相应的行中的实体的 <code>color</code> 属性中。



如果要显示动态文本，请使用特殊类 <code>Table.PlainTextCell</code> 而不是 [Label](gui_Label.html) 组件。它将简化渲染过程并使表格运行更快。



如果 <code>addGeneratedColumn()</code> 方法接收到的参数是未在 XML 描述中声明的列的标识符，则新列的标题将设置如下：



 

<pre><code>carsTable.getColumn(&#34;colour&#34;).setCaption(&#34;Colour&#34;);</code></pre>

 



还可以考虑使用 XML 的 [generator](#gui_Table_column_generator) 属性做更具声明性的设置方案。





[](#gui_Table_requestFocus)
- 
<code>requestFocus()</code> 方法允许将焦点设置在某一行的具体的可编辑字段上。需要两个参数：表示行的实体实例和列的标识符。请求焦点的示例如下：


 

<pre><code>table.requestFocus(item, &#34;count&#34;);</code></pre>

 



[](#gui_Table_scrollTo)
- 
<code>scrollTo()</code> 方法允许将表格滚动到具体行。需要一个参数：表示行的实体实例。


滚动条的示例：



 

<pre><code>table.scrollTo(item);</code></pre>

 



[](#gui_Table_CellClickListener)
- 
如果需要在单元格中显示自定义内容并且在用户单击单元格的时候能收到通知，可以使用 <code>setClickListener()</code> 方法实现这些功能。<code>CellClickListener</code> 接口的实现接收选中实体和列标识符作为参数。这些单元格的内容将被包装在一个 span 元素中，这个 span 元素带有 <code>cuba-table-clickable-cell</code> 样式，可以利用该样式来定义单元格外观。


使用 <code>CellClickListener</code> 的示例：



 

<pre><code>@Inject
private Table&lt;Customer&gt; customersTable;
@Inject
private Notifications notifications;

@Subscribe
protected void onInit(InitEvent event) {
    customersTable.setCellClickListener(&#34;name&#34;, customerCellClickEvent -&gt;
            notifications.create()
                    .withCaption(customerCellClickEvent.getItem().getName())
                    .show());
}</code></pre>

 



[](#gui_Table_setStyleProvider)
- 
<code>setStyleProvider()</code> 方法可以设置表格单元格显示样式。该方法接受 <code>Table.StyleProvider</code> 接口的实现类作为参数。表格的每一行和每个单元分别调用这个接口的 <code>getStyleName()</code> 方法。如果某一行调用该方法，则第一个参数包含该行显示的实体实例，第二个参数为 <code>null</code>。如果单元格调用该方法，则第二个参数包含单元格显示的属性的名称。


设置样式的示例：



 

<pre><code>@Inject
protected Table customersTable;

@Subscribe
protected void onInit(InitEvent event) {
    customersTable.setStyleProvider((customer, property) -&gt; {
        if (property == null) {
        // style for row
        if (hasComplaints(customer)) {
            return &#34;unsatisfied-customer&#34;;
        }
    } else if (property.equals(&#34;grade&#34;)) {
        // style for column &#34;grade&#34;
        switch (customer.getGrade()) {
            case PREMIUM: return &#34;premium-grade&#34;;
            case HIGH: return &#34;high-grade&#34;;
            case MEDIUM: return &#34;medium-grade&#34;;
            default: return null;
        }
    }
        return null;
    });
}</code></pre>

 



然后应该在应用程序主题中设置的单元格和行样式。有关创建主题的详细信息，请参阅 [主题](gui_themes.html)。对于 Web 客户端，新样式在 <code>styles.scss</code> 文件中。在控制器中定义的样式名称，以及表格行和列的前缀标识符构成 CSS 选择器。例如：



 

<pre><code>.v-table-row.unsatisfied-customer {
  font-weight: bold;
}
.v-table-cell-content.premium-grade {
  background-color: red;
}
.v-table-cell-content.high-grade {
  background-color: green;
}
.v-table-cell-content.medium-grade {
  background-color: blue;
}</code></pre>

 



[](#gui_Table_addPrintable)
- 
<code>addPrintable()</code> 当通过 <code>excel</code> [标准操作](standard_actions.html)或直接使用 <code>ExcelExporter</code> 类导出数据到 XLS 文件时，此方法可以给列中数据设置自定义展现。该方法接收的两个参数为列标识符和为列提供的 <code>Table.Printable</code> 接口实现。例如：


 

<pre><code>ordersTable.addPrintable(&#34;customer&#34;, new Table.Printable&lt;Customer, String&gt;() {
    @Override
    public String getValue(Customer customer) {
        return &#34;Name: &#34; + customer.getName;
    }
});</code></pre>

 



<code>Table.Printable</code> 接口的 <code>getValue()</code> 方法应该返回在表格单元格中显示的数据。返回的数据不一定是字符串类型，该方法可以返回其它类型的值，比如数字或日期，它们将在 XLS 文件中以相应的类型展示。



如果生成的列需要在输出到 XLS 时带有格式，则应该使用 <code>addGeneratedColumn()</code> 方法，传递一个 <code>Table.PrintableColumnGenerator</code> 接口的实现作为参数。XLS 文档中单元格的值在这个接口的 <code>getValue()</code> 方法中定义：



 

<pre><code>ordersTable.addGeneratedColumn(&#34;product&#34;, new Table.PrintableColumnGenerator&lt;Order, String&gt;() {
    @Override
    public Component generateCell(Order entity) {
        Label label = uiComponents.create(Label.NAME);
        Product product = order.getProduct();
        label.setValue(product.getName() + &#34;, &#34; + product.getCost());
        return label;
    }

    @Override
    public String getValue(Order entity) {
        Product product = order.getProduct();
        return product.getName() + &#34;, &#34; + product.getCost();
    }
});</code></pre>

 



如果没有以某种方式为生成的列定义 <code>Printable</code> 描述，那么该列将显示相应实体属性的值，如果没有关联的实体属性，则不显示任何内容。



[](#gui_Table_setItemClickAction)
- 
<code>setItemClickAction()</code> 方法能够定义一个双击表格行时将执行的[操作](gui_Action.html)。如果未定义此操作，表格将尝试按以下顺序在其操作列表中查找适当的操作：


  - 
由 <code>shortcut</code> 属性指定给 Enter 键的操作

  - 
<code>edit</code> 操作

  - 
<code>view</code> 操作


如果找到此操作，并且操作具有 <code>enabled=true</code> 属性，则执行该操作。




[](#gui_Table_setEnterPressAction)
- 
<code>setEnterPressAction()</code> 方法可以定义按下 Enter 键时执行的[操作](gui_Action.html)。如果未定义此操作，则表将尝试按以下顺序在其操作列表中查找适当的操作：




  - 
由 <code>setItemClickAction()</code> 方法定义的动作

  - 
由 <code>shortcut</code> 属性指定给 Enter 键的操作

  - 
<code>edit</code> 操作

  - 
<code>view</code> 操作



如果找到此操作，并且操作具有 <code>enabled=true</code> 属性，则执行该操作。





- - -


-  table 的属性 
- 
[align](gui_attributes.html#gui_attr_align) - [aggregatable](#gui_Table_aggregatable) - [aggregationStyle](#gui_Table_aggregationStyle) - [caption](gui_attributes.html#gui_attr_caption) - [captionAsHtml](gui_attributes.html#gui_attr_captionAsHtml) - [columnControlVisible](#gui_Table_columnControlVisible) - [columnHeaderVisible](#gui_Table_columnHeaderVisible) - [contextHelpText](gui_attributes.html#gui_attr_contextHelpText) - [contextHelpTextHtmlEnabled](gui_attributes.html#gui_attr_contextHelpTextHtmlEnabled) - [contextMenuEnabled](#gui_Table_allowPopupMenu) - [css](gui_attributes.html#gui_attr_css) - [dataContainer](gui_attributes.html#gui_attr_dataContainer) - [description](gui_attributes.html#gui_attr_description) - [descriptionAsHtml](gui_attributes.html#gui_attr_descriptionAsHtml) - [editable](#gui_Table_editable) - [enable](gui_attributes.html#gui_attr_enable) - [box.expandRatio](gui_attributes.html#gui_attr_expandRatio) - [height](gui_attributes.html#gui_attr_height) - [id](gui_attributes.html#gui_attr_id) - [multiLineCells](#gui_Table_multiLineCells) - [multiselect](#gui_Table_multiselect) - [presentations](#gui_Table_presentations) - [reorderingAllowed](#gui_Table_reorderingAllowed) - [settingsEnabled](gui_attributes.html#gui_attr_settingsEnabled) - [showSelection](#gui_Table_showSelection) - [sortable](#gui_Table_sortable) - [stylename](#gui_Table_stylename) - [tabIndex](gui_attributes.html#gui_attr_tabIndex) - [textSelectionEnabled](gui_attributes.html#gui_attr_textSelectionEnabled) - [visible](gui_attributes.html#gui_attr_visible) - [width](gui_attributes.html#gui_attr_width)

-  table 的元素 
- 
[actions](#gui_Table_actions) - [buttonsPanel](#gui_Table_buttonsPanel) - [columns](#gui_Table_columns) - [rows](#gui_Table_rows) - [rowsCount](#gui_Table_rowsCount)

- [column](#gui_Table_columns) 元素的属性 
- 
[align](#gui_Table_column_align) - [caption](#gui_Table_column_caption) - [captionProperty](#gui_Table_column_captionProperty) - [collapsed](#gui_Table_column_collapsed) - [dateFormat](gui_DateField.html#gui_DateField_dateFormat) - [editable](#gui_Table_column_editable) - [generator](#gui_Table_column_generator) - [id](#gui_Table_column_id) - [link](#gui_Table_column_link) - [linkInvoke](#gui_Table_column_linkInvoke) - [linkScreen](#gui_Table_column_linkScreen) - [linkScreenOpenType](#gui_Table_column_linkScreenOpenType) - [maxTextLength](#gui_Table_column_maxTextLength) - [optionsDatasource](gui_attributes.html#gui_attr_optionsDatasource) - [resolution](gui_DateField.html#gui_DateField_resolution) - [sortable](#gui_Table_column_sortable) - [visible](gui_attributes.html#gui_attr_visible) - [width](#gui_Table_column_width)

- [column](#gui_Table_columns)的元素 
- 
[aggregation](#gui_Table_column_aggregation) - [formatter](gui_formatter.html)

- [aggregation](#gui_Table_column_aggregation)的属性 
- 
[type](#gui_Table_column_aggregation) - [strategyClass](#gui_Table_column_aggregation_strategyClass)

- [rows](#gui_Table_rows)的属性 
- 
[datasource](gui_attributes.html#gui_attr_datasource)

-  table 的预定义样式 
- 
[borderless](#gui_table_borderless) - [compact](#gui_table_compact) - [no-header](#gui_table_no-header) - [no-horizontal-lines](#gui_table_no-horizontal-lines) - [no-stripes](#gui_table_no-stripes) - [no-vertical-lines](#gui_table_no-vertical-lines) - [small](#gui_table_small)

-  API 
- 
[addGeneratedColumn](#gui_Table_addGeneratedColumn) - [addPrintable](#gui_Table_addPrintable) - [addColumnCollapseListener](#gui_Table_ColumnCollapseListener) - [addSelectionListener](#gui_Table_SelectionListener) - [applySettings](gui_api.html#gui_api_settings) - [generateCell](#gui_Table_generateCell) - [getSelected](#gui_Table_getSelected) - [requestFocus](#gui_Table_requestFocus) - [saveSettings](gui_api.html#gui_api_settings) - [scrollTo](#gui_Table_scrollTo) - [setClickListener](#gui_Table_CellClickListener) - [setEnterPressAction](#gui_Table_setEnterPressAction) - [setItemClickAction](#gui_Table_setItemClickAction) - [setStyleProvider](#gui_Table_setStyleProvider)


- - -


