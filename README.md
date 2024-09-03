# 更新日志
* fix: activemq 解析host:port
* feat: 添加http content type属性
* feat: 解析http响应的content type
* feat: apptypeprocessor支持更多组件
* feat: httpbody解析xml body
* 解析nginxprocessor中的baggage
* 添加spanfault exporter duration
* 修复无效的gap/selfDuration
* 修复spanfault库版本
* 添加spanfault存储属性(gap, selfduration)
* 修复(spanfault, spanaggregation): exporter auth添加username/password
* 修复(nginx): 添加platform.name
* feat(httpbody): 限制http body的值
* feat(httpbodyprocessor): 添加processor
* 添加delaymonitor扩展
* feat(nginxprocessor): 添加nginx processor
* 修复版本0.99 => 0.98
* 添加logstash exporter
* 修复(spanfault): 记录root platform和app cluster
* 测试: 更新测试用例
* 测试: 添加一些测试用例
* feat(apptype): 通过host和port收集server
* 修复(metadata): 不记录log body的值
* feat(apptype): 添加db.connection_string作为server.url
* 升级版本并修复代码
* feat(log): 导出log severity和body以供搜索
* 修复(exception): 添加exception.definition.name
* 修复(metadata): 使用正确的ttl天数
* 修复(metadata): 修复测试
* fix(metadata): 缓存黑名单
* fix(apptype): 使用正确的软件分类
* fix(metadata): 修复导出application structure
* feat(metadata): 只采集不在黑名单中的属性
* fix(metadata): 修复导出
* spanfault: span kind使用正确的字符串
* spanfault: 添加spankind
* spanfault: 记录正常trace
* spanfault: 记录resource/span attributes
* app type processor: 修复bug
* spanfault: 修复nil ref并正确记录FaultKind
* app type processor: 随数据库变更
* exceptionprocessor: 随数据库变更
* app type: 区分应用和客户端的软件类型
* spanfault: 记录root duration/fault span
* exceptionprocessor: 取消锁
* spanaggregation: 使用正确的platform键名
* spanfault: 使用正确的platform键名
* spanfault: 使用groupbytrace简化逻辑
* metadataexport: 控制goroutine数量
* spanaggregation: 使用groupbytrace简化逻辑并提升性能
* metadataexport: 异步导出
* spanaggregation: 改进性能，修复重复保存的bug
* k8sattributesProcessor: 修复cluster info
* metadataexport: 关闭遗漏事务连接
* spanfaultexporter: 修复npe
* metadataexport: 正确关闭事务连接
* spanfaultexporter: 优化队列和查询
* spanfaultexporter: 创建尽可能少的记录
* metadataexport: 导出application structure
* exceptionprocessor: 记录exception.definition.id
* exceptionprocessor: 记录exception.id
* metadataexport: Attributes取值使用单引号
* metadataexport: 添加SpanName
* k8sattributesProcessor添加k8s.cluster.name，取代k8s.cluster.uid
* metadataexport: 记录metrcs/logs属性名
* k8sattributesProcessor添加k8s.service.name
* k8sattributesProcessor添加k8s.cluster.uid
* metadataexport: 记录完整的字段属性名并添加更多字段
* k8sattributesProcessor添加k8s.pod.ip
* metadataexport: 刷新validDate和移除过期key/value
* 添加metadataexporter
* spanmetrics: 根据fault kind统计trace的errorRates
* spanmetrics: 移除delta calls
* 添加span fault
* spanaggregation: 改进监听方案
* spanaggregation: 意外小升级ent
* faultkind: 修复无数据的问题
* spanaggregation: 修复单引号不能保存的问题
* faultkind: 定义SystemFault
* 添加fault kind processor
* 清理span aggregation exporter依赖
* 优化span group命名
* spangroup使用pkg路径，以便在external中使用
* 为exception category processor添加repository
* 自定义模块使用官方版本v0.89.0
* 添加span aggregation exporter
* apptypeprocessor不再监测tomcat/netty
* 添加AppTypeProcessor
* 抽取公共类span group
* add exception processor
* k8sattributesProcessor fix context config
* fix dev
* span metrics add sum and percentile

# OpenTelemetry Collector Contrib

This is a repository for OpenTelemetry Collector components that are not suitable for the  [core repository](https://github.com/open-telemetry/opentelemetry-collector) of the collector.

The official distributions, core and contrib, are available as part of the [opentelemetry-collector-releases](https://github.com/open-telemetry/opentelemetry-collector-releases) repository. Some of the components in this repository are part of the "core" distribution, such as the Jaeger and Prometheus components, but most of the components here are only available as part of the "contrib" distribution. Users of the OpenTelemetry Collector are also encouraged to build their own custom distributions with the [OpenTelemetry Collector Builder](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder), using the components they need from the core repository, the contrib repository, and possibly third-party or internal repositories.

Each component has its own support levels, as defined in the following sections. For each signal that a component supports, there's a stability level, setting the right expectations. It is possible then that a component will be **Stable** for traces but **Alpha** for metrics and **Development** for logs.

## Stability levels

Stability level for components in this repository follow the [definitions](https://github.com/open-telemetry/opentelemetry-collector#stability-levels) from the OpenTelemetry Collector repository.

## Gated features

Some features are hidden behind feature gates before they are part of the main code path for the component. Note that the feature gates themselves might be at different [lifecycle stages](https://github.com/open-telemetry/opentelemetry-collector/tree/main/featuregate#feature-lifecycle).

## Support

Each component is supported either by the community of OpenTelemetry Collector Contrib maintainers, as defined by the GitHub group [@open-telemetry/collector-contrib-maintainer](https://github.com/orgs/open-telemetry/teams/collector-contrib-maintainer), or by specific vendors. See the individual README files for information about the specific components.

The OpenTelemetry Collector Contrib maintainers may at any time downgrade specific components, including vendor-specific ones, if they are deemed unmaintained or if they pose a risk to the repository and/or binary distribution.

Even though the OpenTelemetry Collector Contrib maintainers are ultimately responsible for the components hosted here, actual support will likely be provided by individual contributors, typically a code owner for the specific component.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

Triagers ([@open-telemetry/collector-contrib-triagers](https://github.com/orgs/open-telemetry/teams/collector-contrib-triagers))

- [Benedikt Bongartz](https://github.com/frzifus), Red Hat
- [Jared Tan](https://github.com/JaredTan95), DaoCloud
- [Matt Wear](https://github.com/mwear), Lightstep
- [Murphy Chen](https://github.com/Frapschen), DaoCloud
- Actively seeking contributors to triage issues

Emeritus Triagers:

- [Alolita Sharma](https://github.com/alolita), AWS
- [Gabriel Aszalos](https://github.com/gbbr), DataDog
- [Goutham Veeramachaneni](https://github.com/gouthamve), Grafana
- [Punya Biswal](https://github.com/punya), Google
- [Steve Flanders](https://github.com/flands), Splunk

Approvers ([@open-telemetry/collector-contrib-approvers](https://github.com/orgs/open-telemetry/teams/collector-contrib-approvers)):

- [Anthony Mirabella](https://github.com/Aneurysm9), AWS
- [Antoine Toulme](https://github.com/atoulme), Splunk
- [Bryan Aguilar](https://github.com/bryan-aguilar), AWS
- [Curtis Robert](https://github.com/crobert-1), Splunk
- [David Ashpole](https://github.com/dashpole), Google
- [Yang Song](https://github.com/songy23), DataDog
- [Ziqi Zhao](https://github.com/fatsheep9146), Alibaba

Emeritus Approvers:

- [Przemek Maciolek](https://github.com/pmm-sumo)
- [Ruslan Kovalov](https://github.com/kovrus)

Maintainers ([@open-telemetry/collector-contrib-maintainer](https://github.com/orgs/open-telemetry/teams/collector-contrib-maintainer)):

- [Alex Boten](https://github.com/codeboten), Honeycomb
- [Andrzej Stencel](https://github.com/astencel-sumo), Sumo Logic
- [Bogdan Drutu](https://github.com/bogdandrutu), Snowflake
- [Daniel Jaglowski](https://github.com/djaglowski), observIQ
- [Dmitrii Anoshin](https://github.com/dmitryax), Splunk
- [Evan Bradley](https://github.com/evan-bradley), Dynatrace
- [Juraci Paixão Kröhling](https://github.com/jpkrohling), Grafana Labs
- [Pablo Baeyens](https://github.com/mx-psi), DataDog
- [Sean Marciniak](https://github.com/MovieStoreGuy), Atlassian
- [Tyler Helmuth](https://github.com/TylerHelmuth), Honeycomb

Emeritus Maintainers

- [Tigran Najaryan](https://github.com/tigrannajaryan), Splunk

Learn more about roles in the [community repository](https://github.com/open-telemetry/community/blob/main/community-membership.md).

## PRs and Reviews

When creating a PR please follow the process [described
here](https://github.com/open-telemetry/opentelemetry-collector/blob/main/CONTRIBUTING.md#how-to-structure-prs-to-get-expedient-reviews).

New PRs will be automatically associated with the reviewers based on
[CODEOWNERS](.github/CODEOWNERS). PRs will be also automatically assigned to one of the
maintainers or approvers for facilitation.

The facilitator is responsible for helping the PR author and reviewers to make progress
or if progress cannot be made for closing the PR.

If the reviewers do not have approval rights the facilitator is also responsible
for the official approval that is required for the PR to be merged and if the facilitator
is a maintainer they are responsible for merging the PR as well.

The facilitator is not required to perform a thorough review, but they are encouraged to
enforce Collector best practices and consistency across the codebase and component
behavior. The facilitators will typically rely on codeowner's detailed review of the code
when making the final approval decision.

We recommend maintainers and approvers to keep an eye on the
[project board](https://github.com/orgs/open-telemetry/projects/3). All newly created
PRs are automatically added to this board. (If you don't see the PR on the board you
may need to add it manually by setting the Project field in the PR view).
