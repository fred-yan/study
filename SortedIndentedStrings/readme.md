## 字符串按照等级（缩进级别）进行排序
### 输入字符串如右边所示，输出如右边所示，第一层级进行了排序，第二层级字符串在第一层级下同样进行了排序
|     Original      |       Sorted      |
|-------------------|-------------------|
|Nonmetals          |Alkali Metals      |
|    Hydrogen       |    Lithium        |
|    Carbon         |    Potassium      |
|    Nitrogen       |    Sodium         |
|    Oxygen         |Inner Transitionals|
|Inner Transitionals|    Actinides      |
|    Lanthanides    |        Curium     |
|        Europium   |        Plutonium  |
|        Cerium     |        Uranium    |
|    Actinides      |    Lanthanides    |
|        Uranium    |        Cerium     |
|        Plutonium  |        Europium   |
|        Curium     |Nonmetals          |
|Alkali Metals      |    Carbon         |
|    Lithium        |    Hydrogen       |
|    Sodium         |    Nitrogen       |
|    Potassium      |    Oxygen         |

### 代码运行 python SortedIndentedStrings.py 

参考网址：https://www.kancloud.cn/imdszxs/golang/1509547