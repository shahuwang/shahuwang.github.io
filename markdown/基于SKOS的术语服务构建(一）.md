摘要：

术语服务是指对知识组织系统元数据及其内容进行浏览、查询、应用的各种Web服务系统的统称。在当前网络应用的新形势下，对术语服务的可用性提出了更高的要求，统一的标准，跨平台，跨网络，具备语义等一系列新要求，需要以新的思维和方法去构建术语服务平台。本文描述了面向中文知识组织系统的术语服务的功能设计、系统架构及实施方案，以《汉语主题词表》为例，描述了如何采用标准的SKOS语言以及SKOS扩展来描述中文知识组织系统，并以RESTful Web Service的形式对外提供术语服务。

<!--more-->

关键词：术语服务；SKOS；语义网；Web服务；汉语主题词表

ABSTRACT：Terminology services are a set&nbsp; of&nbsp; web service that present ,retrieve and apply vocabularies with their member terms,concepts and relationships. In the new circumstances, there is a higher calling for a new

way to build terminology service&nbsp; that possess a standard across different organizations and disparate systems ,and also has some semantic with its context.This paper

reports the development of a terminology service system oriented to Chinese knowledge organization systems.The function design,system architecture and key implemen

tation technologies are described in detail.Furthermore ,this paper also describes how to semantically represent traditional Chinese vocabularies(e.g.Chinese Thesaurus) using SKOS language and its customized extension.Finally,this paper descripe some strategy to build terminology service with RESTful Web Service

and its benefit.

KEY WORDS:Terminology service； SKOS；Semantic Web； RESTful Web Service

<a name="_Toc18963"></a><a name="_Toc17538"></a><a name="_Toc17838">目录</a>

[第一章：绪论 1](#_Toc18803 )

[1.1：研究背景及相关概念 1](#_Toc3926 )

[1.1.1：研究背景： 1](#_Toc26396 )

[1.1.2：相关概念： 1](#_Toc597 )

[1.2：研究的目的和意义 3](#_Toc2929 )

[第二章：相关研究述评 4](#_Toc24691 )

[2.1：国外术语服务系统现状 4](#_Toc20057 )

[2.2：国内研究现状 5](#_Toc17427 )

[第三章：受控词表的语义化表示 6](#_Toc21009 )

[3.1：语义网及其相关技术 6](#_Toc4681 )

[3.2：词表的语义化表示 14](#_Toc11258 )

[3.2.1：RDF各实现格式比较 15](#_Toc17669 )

[3.2.2：批量转换的代码实现 21](#_Toc10236 )

[3.3：词表SKOS文档的完整性约束验证 29](#_Toc4991 )

[3.4：词表SKOS文件的存储 31](#_Toc30638 )

[3.5：词表数据的查询 32](#_Toc27435 )

[3.6：小结 32](#_Toc13317 )

[第四章：术语服务的开发 34](#_Toc9524 )

[4.1：Web服务 34](#_Toc3430 )

[4.2：基于Clerezza 框架的术语服务实现 36](#_Toc1218 )

[4.3：基于Jesery开发RESTful Web Service 38](#_Toc9541 )

[第五章：结论 41](#_Toc11429 )

[5.1：总结 41](#_Toc9657 )

[5.1.1：与MARC格式相比的优势 41](#_Toc28437 )

[5.1.2：与XML相比的优势 42](#_Toc3753 )

[5.2：下一步工作 44](#_Toc12451 )

[参考文献： 45](#_Toc11654 )

[致谢： 47](#_Toc8086 )

### 第一章：绪论

#### 1.1：研究背景及相关概念

##### 1.1.1：研究背景：

叙词表、本体、分类表、术语表等各类知识组织系统(Knowledge Organization Systems,简称KOS）在信息资源描述、组织、管理、发现等方面的强大功能，已被图书情报界和相关领域广泛人口。为了让知识组织系统得到更有效、更方便的利用，需要对叙词表等各种知识组织资源进行组织和管理。早期主要是通过纸质印刷品来提供词表的使用，后来则提供了光盘的方式。自1996年起，国外出现了以电子格式发布的在线词表列表,如英属哥伦比亚大学图书情报学院的词表索引和HILT Resource List。但是他们都未得到持久的扩展和维护，使得使用起来越来越不方便。自1998年以来，图书情报学界开始了对网络知识组织系统（Networked Knowledge Organization Systems,简称NKOS）的研究，即探索在网络环境下，让知识组织系统成为支持信息资源检索和描述的网络交换式服务。术语服务通过Web应用程序接口（Application Programming Interfaces，简称APIs）支持其他系统对知识组织系统资源的访问和调用，为网络环境下的编目、元数据创建、信息检索、知识组织等各种应用提供强大的术语支持，是传统知识组织系统在网络环境下进行应用的新途径。

本文主要探讨以《汉语主题词表》为例，探讨面向中文网络知识组织系统术语服务的构建，使中文知识组织资源能够发挥其强大术语支持功能和语义功能。

##### 1.1.2：相关概念：

（1）术语

在英文里，术语的单词是Term ，对应的是术语学，即Terminology，比较好的定义是龚益在其《术语、术语学和术语标准》一文做出的解释：术语（term）是在特定学科领域用来表示概念的称谓的集合，或者说，是通过语音或文字来表达或限定科学概念的约定性语言符号。在我国，人们习惯称其为“名词”。术语是传播知识、技能，进行社会文化、经济交流等不可缺少的重要工具。作为科学发展和交流的载体，术语是科学研究的成果，是人类进步历程中知识语言的结晶。从某种意义上说，术语工作的进展和水平，直接反映了全社会知识积累和科学进步的程度。术语和文化，如影之随形，须臾不离。不同的文化要用不同的术语来说明，吸收外来文化，同时必须吸收外来术语。随着社会的发展进步，新概念大量涌现，必须用科学的方法定义、指称这些概念。所谓概念，是客体的抽象，在专门语言中用称谓表示，并用定义描述。研究概念、概念定义和概念命名基本规律的边缘学科，即是术语学，在20世纪30年代初期正式创立。

（2）受控词表

A controlled vocabulary is a way to insert an interpretive layer of semantics between the term entered by the user and the underlying database to better represent the original intention of the terms of the user ，简单的说，受控词表是对术语进行组织、标记，以便用户可更方便进行检索。

（3）术语服务

所谓术语服务是指对知识组织系统元数据及其内容进行浏览、查询和应用的各种Web服务的统称。术语服务通过Web应用程序接口（Application Programming Interfaces，简称APIs）接口支持机器对知识组织系统的访问和应用，能够在网络环境下为编目、元数据创建、信息检索、知识组织等各种应用提供强大的术语支持，是传统的知识组织系统在网络环境下进行应用的新途径。

通过上述几个概念的简介，虽然不能让读者有很好的认知，但至少可以让读者知道术语服务的大概是什么意思。经过多年的发展，现在已经有很多非常优秀的受控词表，比如汉语主题词表，中国图书馆分类号这一类的受控词表。然而，这类受控词表建立的时间都比较早，同时，由于受控词表所处的环境因素，导致受控词表的应用技术上，很难跟上时代发展的步伐。但是，作为一种基础服务，受控词表又异常重要。研究如何使得术语服务能跟上时代步伐，甚至超前一些，使得受控词表能得到更好的利用，便是本课题的主题。

#### 1.2：研究的目的和意义

经过多年的发展，现在已经有很多非常优秀的受控词表，比如《汉语主题词表》，《中国图书馆分类法》，然而这类受控词表建立的时间都比较早，同时由于受控词表所处的环境因素，导致受控词表在应用技术上很难跟上时代发展的步伐。但是作为一种基础服务，受控词表又异常重要。研究如何使得术语服务能跟上时代步伐，甚至超前一些，使得受控词表能得到更好的利用，便是本课题的主题。

本文的目的是针对传统的收控词表构建一个术语服务。在本研究中， 选择《汉语主题词表》中的部分数据作为术语样本数据，选择采用SKOS语言对词表中的术语进行语义化描述,重点介绍了如何在工程上实现将《汉语主题词表》转换为SKOS格式描述，以及存储、查询等基于Semantic Web技术的解决方案。同时，对术语服务的发布形式进行了探究。本文对基于SKOS的术语服务全技术做主要说明，为术语服务的实践应用探明了方向。