package api

type Type int

const (
	// login check
	// 登录检查
	IS_LOGIN Type = 0

	// self info
	// 获取个人信息
	GET_SELF_INFO Type = 1

	// send message
	// 发送文本
	MSG_SEND_TEXT Type = 2
	// 发送群艾特
	MSG_SEND_AT Type = 3
	// 分享好友名片
	MSG_SEND_CARD Type = 4
	// 发送图片
	MSG_SEND_IMAGE Type = 5
	// 发送文件
	MSG_SEND_FILE Type = 6
	// 发送xml文章
	MSG_SEND_ARTICLE Type = 7
	// 发送小程序
	MSG_SEND_APP Type = 8

	// receive message

	// 开启接收消息HOOK，只支持socket监听
	MSG_START_HOOK Type = 9
	// 关闭接收消息HOOK
	MSG_STOP_HOOK Type = 10
	// 开启图片消息HOOK
	MSG_START_IMAGE_HOOK Type = 11
	// 关闭图片消息HOOK
	MSG_STOP_IMAGE_HOOK Type = 12
	// 开启语音消息HOOK
	MSG_START_VOICE_HOOK Type = 13
	// 关闭语音消息HOOK
	MSG_STOP_VOICE_HOOK Type = 14

	// contact
	// 获取联系人列表
	CONTACT_GET_LIST Type = 15
	// 检查是否被好友删除
	CONTACT_CHECK_STATUS Type = 16
	// 删除好友
	CONTACT_DEL Type = 17
	// 从内存中获取好友信息
	CONTACT_SEARCH_BY_CACHE Type = 18
	// 网络搜索用户信息
	CONTACT_SEARCH_BY_NET Type = 19
	// wxid加好友
	CONTACT_ADD_BY_WXID Type = 20
	// v3数据加好友
	CONTACT_ADD_BY_V3 Type = 21
	// 关注公众号
	CONTACT_ADD_BY_PUBLIC_ID Type = 22
	// 通过好友请求
	CONTACT_VERIFY_APPLY Type = 23
	// 修改备注
	CONTACT_EDIT_REMARK Type = 24

	// chatroom
	// 获取群成员列表
	CHATROOM_GET_MEMBER_LIST Type = 25
	// 获取指定群成员昵称
	CHATROOM_GET_MEMBER_NICKNAME Type = 26
	// 删除群成员
	CHATROOM_DEL_MEMBER Type = 27
	// 添加群成员
	CHATROOM_ADD_MEMBER Type = 28
	// 设置群公告
	CHATROOM_SET_ANNOUNCEMENT Type = 29
	// 设置群聊名称
	CHATROOM_SET_CHATROOM_NAME Type = 30
	// 设置群内个人昵称
	CHATROOM_SET_SELF_NICKNAME Type = 31

	// database
	// 获取数据库句柄
	DATABASE_GET_HANDLES Type = 32
	// 备份数据库
	DATABASE_BACKUP Type = 33
	// 数据库查询
	DATABASE_QUERY Type = 34

	// version
	// 修改微信版本号
	SET_VERSION Type = 35

	// log
	// 开启日志信息HOOK
	LOG_START_HOOK Type = 36
	// 关闭日志信息HOOK
	LOG_STOP_HOOK Type = 37

	// browser
	// 打开微信内置浏览器
	BROWSER_OPEN_WITH_URL Type = 38
	// 获取公众号历史消息
	GET_PUBLIC_MSG Type = 39

	// 转发消息
	MSG_FORWARD_MESSAGE Type = 40
	// 获取二维码
	GET_QRCODE_IMAGE Type = 41
	// 获取A8Key
	GET_A8KEY Type = 42
	// 发送xml消息
	MSG_SEND_XML Type = 43
	// 退出登录
	LOGOUT Type = 44
	// 收款
	GET_TRANSFER Type = 45
)
