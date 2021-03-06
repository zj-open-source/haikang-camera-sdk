# Golang海康工业相机SDK二次开发 海康sdk二次开发(golang版)

[海康工业相机SDK下载](https://www.hikrobotics.com/cn/machinevision/service/download?module=0)

+ SDK版本:基于海康SDK中心MVS V2.1.0(Linux)版本
+ 操作系统 Ubuntu18.04
+ 环境 golang16.7
+ 硬件设备 海康线阵相机,面阵相机,旋转编码器(触发相机拍照--非必须)

# 调试准备

+ 下载好MVS到Ubuntu系统中,打开MVS软件,将网卡设置在同一个局域网的不同网段中(我们机器是四网卡),确保相机能再MVS软件中被枚举出来,才能继续运行程序.
+ 若MVS中能够枚举出相机,但是运行程序失败,请先检查环境变量是否正确:$LD_LIBRARY_PATH=/opt/lib/64/:/opt/lib/32/,否则请执行source /opt/MVS/bin/set_env_path.sh

# 调试

+ 项目根目录的grabImage_test.go中的两个测试方法可以直接运行;
+ 如果程序正常运行,但是没有图片回调或生成,请将相机触发模式设置为关闭(即自动触发取图);

+ 项目路径__test__/hikvision/main.go程序可以正常运行

# 代码简单说明

+ 官方Demo只有C++和Python，该项目仅仅是海康工业相机SDK的GO版本,根据Python的demo翻译而来
+ 该demo不区分线阵和面阵相机,均适用两种相机
+ MacOS的mvs_sdk在/Library/MVS_SDK下面
+ demo在camera-sdk.go文件中定义了一个公共的结构体，海康相机的SDK集成在camera_sdk_hikvision.go中实现，还可以扩展除海康之外的其他相机SDK
+ 作者水平有限,有些处理方式不是最好的,欢迎朋友指正

# 统计

![Alt](https://repobeats.axiom.co/api/embed/5eaf96a7329e3c8622dc77bd6065c6356306bd48.svg "Repobeats analytics image")

# TODO list

+ camera_sdk_hikvision.go中不同的图像格式,buffSize是不一样的,抽象出来,并且格式是mvs的客户端设置的;
+ 海康工业相机的软触发模块还缺少了,找个时间把软触发这一块给补上.软触发的demo(Mac:/Library/MVS_SDK/Samples/Trigger_Image,Ubuntu:
  /opt/mvs/Samples/Trigger_Image)