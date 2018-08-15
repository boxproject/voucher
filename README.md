![architecture](https://github.com/boxproject/voucher/blob/master/banner.png)

# PRODUCT INFORMATION

**Product name**: BOX

**Description**: Enterprise-grade security solution for digital assets custody, using a cryptographically secure offline network.

**Version**: 1.0

[中文版说明](https://github.com/boxproject/voucher/blob/master/README_CN.md)

# BRIEF MARKET INTRODUCTION

Banking and financial organizations are paying more attention to cryptocurrency, in particular hedge funds. The financing scale of VC investing in digital assets has shoot up from USD 2 million in 2012 to USD 3.4 billion, representing an increase by 1672.5 times in 6 years. The accumulative financing scale of VC digital assets investment climbed from 3 in 2012 to 182 in 2018. Institutional investors and businesses are now entering a crypto market becoming mature because their interest in funds dealing with blockchain technology and cryptocurrencies has grown steadily in recent months.

As a consequence, it is expected that funding for cryptocurrency funds may grow by 1,500% by Q4 2018 as digital assets are now outperforming traditional financial assets and various banking and financial organizations invest massively in blockchain and altcoins projects. International payment systems providers are working on developing digital transactions. This will allow and push the world’s largest financial and banking organizations to invest big figures, further expanding the market for digital coins and impacting their value.

# REGULATORY CONSTRAINTS

For years the digital assets market has go on unregulated, and recent actions that some government have taken to 
Steps that some governments have taken to prevail the risks of the crypto-market seem to have increased the confidence of institutional investors. Institutional investors prefer investing in markets containing regulatory frameworks.

The U.S. Securities and Exchange Commission (SEC) have started implementing securities laws to regulate the crypto-market and protect customers. This has created an environment where this investor class can access to capital-intensive, long-term investment opportunities.

# PROBLEMS FACED BY INSTITUTIONS

When it comes to safekeeping, the investor class is highly relying either on personal wallets or cold wallets, causing many unexpected problems. While Mobile wallets are more practical and easier to use than other crypto wallet types thanks to its features of accepting and sending payments on the fly, phones and laptops are insecure devices that can be stolen, maliciously compromised (malware, keyloggers, etc.) or rooted. When it comes to Cold wallets, they require more effort to move cryptocurrencies around, also more technical knowledge.

The security of such private wallets cannot be fully guaranteed against malicious issues. As for the compliance side, the use of personal wallets by organizations for managing digital currencies is not compatible with any local standardized financial management process, leading to confusions and mistakes with accounting documents filing.

For security reasons, organizations have preferred to gradually shift from hot wallets to cold wallets. Due to its trawl characteristics and problems explain above, its security level remains higher than software wallets. In order to ensure the absolute safety of the cold wallet, they are usually kept in banks’ safe, which causes great inconvenience for making multiple transfers daily. In addition, operating transfers with a cold wallet can be complicated for non-professionals, resulting in inefficient use of the cold wallet or worst, mistakes. If the cold wallet contains an embedded multi-signature technology providing co-management of the private key, it also has some defects because the main chain wallets do not support this technology.

Under the circumstances that institutions usually hold different cryptocurrencies, multi-signature poses some obstacles to manage the institutional digital assets. In addition, the introduction of multi-signature technologies in different main chains results in non-portability of digital asset management. More critical, there are loopholes in some of the main chains using multiple signatures, such as Ethereum's PARITY event. In conclusion, the market is in need for a comprehensive solution which protects efficiently the company’s digital assets.

# SOLUTION WITH ISOMETRIC DEVICE MOCKUP

BOX provide an enterprise-grade security solution for digital assets custody, using cryptographically secure offline network including flow of approvals, private blockchain technology and communication security. BOX achieve integration of technologies and fundamentally solve the industry security issues such as the theft of private keys and the falsification of directives.

![architecture](https://github.com/boxproject/voucher/blob/master/architecture.png)

![example](https://github.com/boxproject/voucher/blob/master/process.png)

**Owning the private key of an account gives full access to the fund.** The dynamic password provides shared authority, a one-click activation and ensures the security of private keys. The BOX system uses a single private key to manage all cryptocurrencies. Theoretically, all public chains that support the ECDSA elliptic curve algorithm can be controlled with the private key. At this point, the BOX system is more convenient than multi-signature. Meanwhile, BOX uses a multi-person multi-password method to automatically generate a private key by using an algorithm in a signature machine, and then generate a public key from the private key. The partners who have the highest authority only have a part of the dynamic password which provides him with a shared governance on the private key.

In terms of storage, we put the private key in the memory of the signature machine, without any persistent storage, thus making it extremely difficult to be captured. We lock the private key in the memory to prevent bypass attacks. In the event of a power outage, the BOX’s signature machine will automatically shut down the memory and the private key will disappear.

Therefore, even if the signature machine is streaking, the chance of obtaining the private key from the BOX system is almost zero. The partner with the highest authority can instantly restore the original dynamic password by putting it in banks’ safe, in order to prevent a partner from accidentally failing to perform duties. Unlike cold wallets, there is no need to move this backup frequently. Only when a partner has an accident will he vote via the board to decide whether to enable password backup.

The custom approval flow template uses the features of the blockchain that cannot be tampered with to be stored in the private chain. The template of the approval flow is defined by the enterprise itself. The content mainly includes the hierarchy, the initiation (approval), the minimum number of employees at each level, and the employee's public key (address). As a result, the hash value of the custom template and the template on the private chain both ensure that the approval flow cannot be modified. The private key APP will confirm its validity. When an employee initiates an approval flow, if the employee's private key and the address corresponding to the private key are matched, the approval process is matched with the approval flow template stored on the private chain through the associated program on the private chain. If it is in full compliance with the approval flow template, and then through the proxy (private key app interface) flows to the signing machine, before the transfer of the private key in the signing machine, and the hash of the approval flow template stored in the public chain is matched twice (currently The secondary matching of the public chain only supports Ethereum). After ensuring that there is no mistake, the private key in the signature machine will be transferred for transfer. In addition, BOX provides a unified public account for each company, so that the company's assets can be managed under one account for effective management. All digital assets will be traded through the account, preventing the case that public and private accounts are not separated. The approval flow also provides the basis for the audit, the company's managers can also clearly understand the company's assets through these records, and conduct a corresponding analysis.

On the hardware side, deploying a BOX system requires at least 3(2n+1) cloud servers. Each cloud server acts as a node and builds a private chain. An Apple MACBOOK as a signing machine, because IOS is more secure than Windows. Several iPhones are needed to load the private key APP & employee APP.

**One-stop integrated solution, BOX system will safeguard the security of investment firms, crypto-exchange platforms and other digital assets.** At present, the most suitable companies for BOX investment firms interested in Blockchain, Blockchain companies with audit risk control/compliance requirements, and trading platforms. Blockchain investment firms often transfer funds and receive payments frequently. It is inconvenient to use cold wallets. Personal wallets are not suitable for institutions.

**The BOX code has been uploaded to GitHub, the largest technology open source community in the world, and it is necessary to build a healthier and safer industry environment together with the you.** Any individual or enterprise can use and deploy the system free of charge. In order to stimulate the first contributors of the BOX 0.1.0 version, the BOX team launched the “BOX Super Partner” program. Up to now, more than 30 organizations have signed letters of intent with the BOX Foundation. In the future, the BOX team will focus on community building and system scalability, and work with many organizations to build a healthier and safer industry environment. BOX's open source code on github has multiple repositories providing a complete set of deployable solutions. Including agent - private key APP management server, box-Authorizer - private key APP client, boxguard - signature machine daemon, voucher - access layer, companion - private chain side companion program, box-Staff-Manager - employee APP client End, box-appServer - Employee APP Server.

## Main functionalites

There are 5 main functionalites for this voucher progam:

1. Communicate with proxy server on private chain
2. Upload realtime status
3. Receive request from private key app, complete offline signature and publish smart contract etc.
4. Submitted data will only be confirmed when transctions or approval flows initialized from private blockchain, signed and verified using RSA.
5. Monitor event log of Ethereum, confirm approval flows, topup and withdrawl and transfer result to proxy server.

**How to use:**

1. Initialization. 
Use the following command for the first time run
```sh
make build
cp config.toml.example  config.toml
cp log.xml.example  log.xml
```
Otherwise, use this command
```sh
make rebuild
```
Note: build command will clear all data.

2. Update parameters (address to associated program, port, public key and certificate) in config.json respectively. 
3. Update rpc api address to config.json
4. Run program


**Run command in cli：**

```bash
➜ ./voucher
```

# Legal Reminder

Exporting/importing and/or use of strong cryptography software, providing cryptography hooks, or even just communicating technical details about cryptography software is illegal in some parts of the world. If you import this software to your country, re-distribute it from there or even just email technical suggestions or provide source patches to the authors or other people you are strongly advised to pay close attention to any laws or regulations which apply to you. The authors of this software are not liable for any violations you make - it is your responsibility to be aware of and comply with any laws or regulations which apply to you.

We would like to hear critics and feedbacks from you as it will help us to improve this open source project.
