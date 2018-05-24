# BOX Product Introduction

**The investment made by institutions into digital assets has boomed in recent years.** The financing scale of VC investing in digital assets has soared from USD 2 million in 2012 to USD 3.4 billion, representing an increase by 1672.5 times in 6 years. The accumulative financing scale of VC digital assets investment surged from 3 in 2012 to 182 in 2018. According to recent trends for asset digitalization, it is expected that more organizations will aim for it in the near future.

**Institutional digital asset management is highly relying on personal wallets, causing many unexpected problems.** At present, digital asset management of enterprises also relies on personal wallets such as Imtoken. The security of such private wallets cannot be fully guaranteed, let alone effective prevention of theft. In terms of compliance, the use of personal wallets by organizations for managing digital currencies could not achieve a standardized financial management process, leading to confusions and mistakes with accounting documents filing. 

**For security reasons, organizations have gradually shifted from online softwares to cold wallets.** Due to the trawl characteristics of the cold wallet, its security level is higher than software wallets. In order to ensure the absolute safety of the cold wallet, cold wallets are usually kept in banks’ safe, which causes great inconvenience for making multiple transfers daily. In addition, operating transfers with a cold wallet can be complicated for non-professionals, resulting in inefficient use of the cold wallet or worst, mistakes.

The embedded multi-signature technology of the wallet provides co-management of the private key, but it also has some defects because the main chain wallets do not support this technology.

**Under the circumstances that institutions usually hold different cryptocurrencies, multi-signature poses some obstacles to manage the institutional digital assets.** In addition, the introduction of multi-signature technologies in different main chains results in non-portability of digital asset management. More critical, there are loopholes in some of the main chains using multiple signatures, such as Ethereum's PARITY event. In conclusion, the market is in need for a comprehensive solution which protects efficiently the company’s digital assets.

**BOX provide a secured one-stop solution for managing enterprise digital assets, based on blockchain technology, cryptography, and communication security.** BOX achieve integration of technologies and fundamentally solve the industry security issues such as the theft of private keys and the falsification of directives. The security and efficiency of the entire BOX system is mainly achieved through private keys, a flow of approvals, and crypted communication.

![整体架构图](https://s3-ap-southeast-1.amazonaws.com/s3.box.images01/QQ20180518-141013%402x.png)

![审批流程示意](https://s3-ap-southeast-1.amazonaws.com/s3.box.images01/QQ20180518-140745%402x.png)

**Owning the private key of an account gives full access to the fund.** The dynamic password provides shared authority, a one-click activation and ensures the security of private keys. The BOX system uses a single private key to manage all cryptocurrencies. Theoretically, all public chains that support the ECDSA elliptic curve algorithm can be controlled with the private key. At this point, the BOX system is more convenient than multi-signature. Meanwhile, BOX uses a multi-person multi-password method to automatically generate a private key by using an algorithm in a signature machine, and then generate a public key from the private key. The partners who have the highest authority only have a part of the dynamic password which provides him with a shared governance on the private key.

In terms of storage, we put the private key in the memory of the signature machine, without any persistent storage, thus making it extremely difficult to be captured. We lock the private key in the memory to prevent bypass attacks. In the event of a power outage, the BOX’s signature machine will automatically shut down the memory and the private key will disappear.

Therefore, even if the signature machine is streaking, the chance of obtaining the private key from the BOX system is almost zero. The partner with the highest authority can instantly restore the original dynamic password by putting it in banks’ safe, in order to prevent a partner from accidentally failing to perform duties. Unlike cold wallets, there is no need to move this backup frequently. Only when a partner has an accident will he vote via the board to decide whether to enable password backup.

The custom approval flow template uses the features of the blockchain that cannot be tampered with to be stored in the private chain. The template of the approval flow is defined by the enterprise itself. The content mainly includes the hierarchy, the initiation (approval), the minimum number of employees at each level, and the employee's public key (address). As a result, the hash value of the custom template and the template on the private chain both ensure that the approval flow cannot be modified. The private key APP will confirm its validity. When an employee initiates an approval flow, if the employee's private key and the address corresponding to the private key are matched, the approval process is matched with the approval flow template stored on the private chain through the associated program on the private chain. If it is in full compliance with the approval flow template, and then through the proxy (private key app interface) flows to the signing machine, before the transfer of the private key in the signing machine, and the hash of the approval flow template stored in the public chain is matched twice (currently The secondary matching of the public chain only supports Ethereum). After ensuring that there is no mistake, the private key in the signature machine will be transferred for transfer. In addition, BOX provides a unified public account for each company, so that the company's assets can be managed under one account for effective management. All digital assets will be traded through the account, preventing the case that public and private accounts are not separated. The approval flow also provides the basis for the audit, the company's managers can also clearly understand the company's assets through these records, and conduct a corresponding analysis.

On the hardware side, deploying a BOX system requires at least 3(2n+1) cloud servers. Each cloud server acts as a node and builds a private chain. An Apple MACBOOK as a signing machine, because IOS is more secure than Windows. Several iPhones are needed to load the private key APP & employee APP.

**One-stop integrated solution, BOX system will safeguard the security of investment firms, crypto-exchange platforms and other digital assets.** At present, the most suitable companies for BOX investment firms interested in Blockchain, Blockchain companies with audit risk control/compliance requirements, and trading platforms. Blockchain investment firms often transfer funds and receive payments frequently. It is inconvenient to use cold wallets. Personal wallets are not suitable for institutions.

**The BOX code has been uploaded to GitHub, the largest technology open source community in the world, and it is necessary to build a healthier and safer industry environment together with the you.** Any individual or enterprise can use and deploy the system free of charge. In order to stimulate the first contributors of the BOX 0.1.0 version, the BOX team launched the “BOX Super Partner” program. Up to now, more than 30 organizations have signed letters of intent with the BOX Foundation. In the future, the BOX team will focus on community building and system scalability, and work with many organizations to build a healthier and safer industry environment. BOX's open source code on github has multiple repositories providing a complete set of deployable solutions. Including agent - private key APP management server, box-Authorizer - private key APP client, boxguard - signature machine daemon, voucher - access layer, companion - private chain side companion program, box-Staff-Manager - employee APP client End, box-appServer - Employee APP Server.

# BOX产品介绍

**机构掀起数字资产投资狂潮。** 全球 VC 数字资产投资规模由2012年2百万美元飙升至目前34亿美元，期间增幅达1672.5倍。其中，VC数字资产投资累计融资规模由2012年3起激增至2018年182起。在资产数字化的大浪潮下，预计未来将有更多的机构参与进来。

**机构数字资产管理高度依赖个人钱包，引发诸多问题。** 目前企业数字资产管理还依赖于Imtoken等个人钱包，该类钱包目前连私钥的安全都无法完全保证，更别提做到有效防范内鬼。合规方面，机构使用个人钱包管理数字货币无法实现规范化的财务管理流程，导致账目混乱。此外，在缺乏指导与监督的情况下，数字货币打错地址，资金无法追回的意外时有发生。

出于数字货币安全层面的考虑，机构数字货币管理逐渐向冷钱包转移。由于冷钱包的拖网特性，其安全等级要高于热钱包。为了确保冷钱包的绝对安全，通常情况下都要将冷钱包保存于银行保险库 ，这为机构的高频率打款转账造成了极大的不便。此外，冷钱包的操作方法对非专业人士而言较为复杂，致使冷钱包使用效率低下。

**钱包嵌入多重签名技术实现私钥多人共管，但便捷性等方面仍存缺陷。** 钱包嵌套多重签名技术一定程度实现了机构多人共管私钥的目的。但多重签名目前依旧存在如下问题：首先，主链需要支持多重签名技术，主链钱包也需支持多重签名，然而现阶段许多主链并不支持多重签名技术；在机构普遍持有多种数字货币的情况下，多重签名给机构数字资产的管理造成一定阻碍。此外，不同主链导入多重签名技术也不同，造成了数字资产管理的极大的不便携性。更为致命的是，有些主链的多重签名会有漏洞，比如以太坊的PARITY事件。综合看，目前市面亟需一款专门针对企业级的数字资产管理的综合解决方案，为企业数字资产管理保驾护航。

**BOX的横空出世，为企业数字资产安全的管理提供了一套综合解决方案。** BOX(企业通证保险箱)是一套基于区块链技术的数字资产安全解决方案。BOX以区块链、密码学、通信安全等领域的公理性技术为依托，实现技术间的无缝衔接，从根本上解决了私钥被盗取和指令被篡改的行业痛疾。整套BOX系统的安全高效主要通过私钥、审批流、以及通信安全来实现。

**一把私钥管理几乎所有数字货币，动态口令方式实现多人共管，一键启停确保私钥安全。** BOX系统实现了一把私钥管理所有数字货币的目标，理论上所有支持ECDSA的椭圆形曲线算法的公链均可以用一把私钥来管理，在这一点上，BOX系统较多重签名在便携性方面有了巨大的提升。同时，BOX采用多人多口令方式，在签名机内通过算法将口令动态生成私钥，再用私钥去生成公钥。拥有最高权限的合伙人均只掌握动态口令的一部分，实现了私钥多人共管的目的。在存储方面，我们将私钥放在了签名机内存里，不做任何持久化存储，内存里的私钥极难被捕捉。我们又将私钥内存位两端锁住，防止旁道攻击。一旦发生断电或启动等情况，BOX签名机的守护进程将自动清内存关机，这样私钥就消失了。因此即使在签名机裸奔的情况下，从BOX系统得到私钥的机会也几乎为零。拥有最高权限的合伙人可以通过依次输入口令的方式即刻再次恢复原有私钥。此外，为了防止某个合伙人出现意外无法履行职责，我们推荐在设置口令时，通过打印口令备份封存到银行保险箱，并且将其中一份备份交给三方保存。与冷钱包不同的是，平时不需要去动这个备份，仅当合伙人发生意外时才会经由董事会投票决定是否启用口令备份。

![整体架构图](https://s3-ap-southeast-1.amazonaws.com/s3.box.images01/QQ20180518-141013%402x.png)

![审批流程示意](https://s3-ap-southeast-1.amazonaws.com/s3.box.images01/QQ20180518-140745%402x.png)

**自定义审批流模版，利用区块链不可篡改的特性存入私链中。统一对公账户，高效管理数字资产。** 审批流的模版由企业自主定义，内容主要包括审批层级、发起（审批）人、每一层级最少通过人数、员工的公钥（地址）。将自定义好的模版以及模版的哈希值保存在私链上，以确保审批流的不可篡改性，再由私钥APP确认其有效性。当员工发起审批流时，首先匹对员工私钥以及私钥所对应的地址，确保无误后，通过私链上的伴生程序，将审批流程与保存在私链上的审批流模版一一匹对，若完全符合审批流模版，再通过代理（私钥app接口）流向签名机，在调动签名机里的私钥前，与保存在公链上的审批流模版哈希值进行二次匹配（目前公链二次匹配只支持以太坊），确保无误后，将调动签名机里的私钥进行打款转账。除此之外，BOX为每一个企业提供了统一的对公账户，这样能把企业的资产放在一个账户下进行有效的管理，所有的数字资产的交易都将通过该账户进行交易，防止各种公账私账不分的情况的发生。多数字资产集中管理和数字资产交易的明细这两点，主要体现了BOX系统能够实现多种数字资产的统一管理，并为这些数字资产的交易提供明细记录，这样可以帮助企业进行的账务记录，并为审计提供依据，企业的管理者也可以通过这些记录非常清楚的了解企业资产情况，并进行相应的分析。这些都对于数字资产的有效管理非常重要。

硬件方面，部署一套BOX系统需要至少 3（2n+1）台云服务器，每一个云服务器作为一个节点，构建一条私链。一台苹果 MACBOOK 作为签名机，采用苹果电脑作为签名机是考虑到苹果系统较Windows更为安全。同时需要若干台iphone，装载私钥APP & 员工APP。

**一站式综合解决方案，BOX系统将为投资机构、交易所等数字资产安全保驾护航。** 目前BOX最适合的企业是区块链投资机构、有区块链审计风控需求的企业、交易平台等。区块链投资机构平时会频繁转账、收款，用冷钱包不方便，个人钱包又不适用于多位合伙人企业进行清晰的资产划分，BOX则是多人共管一把私钥，可以很好的解决这一难题，操作起来也比冷钱包便利很多。总体来看，BOX系统通过一系列的优化流程，大幅提高可用性及便捷性，同时达到与冷钱包相当的安全程度。

**BOX代码已上传至全球最大的技术开源社区GitHub，待与诸君共同构建更健康、安全的行业大环境。** 任何个体、企业均可无偿使用、部署该系统。为激励 BOX 0.1.0 版本的首批贡献者，BOX团队启动了“BOX超级合伙人“计划，截至目前已有30余家机构与BOX基金会签署了合作意向书。未来BOX团队将着重社群建设以及系统的可扩展性，与众多机构共同构建更健康、安全的行业大环境。BOX在github上的开源代码有多个仓库，提供了整套可部署的解决方案。包括 agent - 私钥APP管理服务端，box-Authorizer - 私钥APP客户端，boxguard - 签名机守护程序，voucher - 接入层，companion - 私链端伴生程序，box-Staff-Manager - 员工APP客户端， box-appServer - 员工APP服务器端。

## 核心功能

本程序主要功能点有五个：

1. 与代理服务通信
2. 实时上报签名机状态信息
3. 接受私钥app的请求，完成离线签名，发布智能合约等
4. 只接受经过私链确认的审批流以及转账，并进行RSA验签，才认可本次提交的数据正确性。
5. 监控以太坊公链event log，确认审批流操作，以及充值提现记录等，并将结果通知给代理


**使用步骤：**

1. 初始化。首次使用make build命令，再次make rebuild命令，切记，build命令将清除所有数据；
2. 将连接代理的地址、端口以及ssl公钥以及证书写入到config.json配置文件中对应的参数中
3. 将eth公链、比特币公链分别写入到config.json配置文件中对应的参数中
4. 启动本程序


**命令行用法：**

```bash
➜ ./voucher
```


**功能说明：**

＊ 故障恢复

＊ 根据配置设定出块数据确认，以防止因分叉导致确认数目出错。

＊ 审批流创建及确认监控

＊ 充值、转账审批监控
