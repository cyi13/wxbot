package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Api struct {
	host    string
	timeout time.Duration
}

type Config struct{}

func New(host string) *Api {
	if !strings.HasPrefix(host, "http") {
		host = "http://" + host
	}
	return &Api{
		host:    host + "/api/?type=%d",
		timeout: 5 * time.Second,
	}
}

type IsLoginResult struct {
	IsLogin int    `json:"is_login"`
	Result  string `json:"result"`
}

func (a *Api) RequestIsLogin() (*IsLoginResult, error) {
	res := &IsLoginResult{}
	if err := a.sendApi(IS_LOGIN, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type GetSelfInfoResult struct {
	Data struct {
		PhoneNumber   string `json:"PhoneNumber"`
		Sex           string `json:"Sex"`
		Uin           int    `json:"uin"`
		WxBigAvatar   string `json:"wxBigAvatar"`
		WxCity        string `json:"wxCity"`
		WxFilePath    string `json:"wxFilePath"`
		WxID          string `json:"wxId"`
		WxNation      string `json:"wxNation"`
		WxNickName    string `json:"wxNickName"`
		WxNumber      string `json:"wxNumber"`
		WxProvince    string `json:"wxProvince"`
		WxSignature   string `json:"wxSignature"`
		WxSmallAvatar string `json:"wxSmallAvatar"`
	} `json:"data"`
	Result string `json:"result"`
}

// 获取个人信息
func (a *Api) RequestGetSelfInfo() (*GetSelfInfoResult, error) {
	res := &GetSelfInfoResult{}
	if err := a.sendApi(GET_SELF_INFO, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type MsgSendTextData struct {
	Wxid string `json:"wxid"`
	Msg  string `json:"msg"`
}

type MsgSendTextResult struct {
	Msg    int    `json:"msg"`
	Result string `json:"result"`
}

// 发送文本
func (a *Api) RequestMsgSendText(data *MsgSendTextData) (*MsgSendTextResult, error) {
	res := &MsgSendTextResult{}
	if err := a.sendApi(MSG_SEND_TEXT, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type MsgSendAtData struct {
	ChatroomId   string `json:"chatroom_id"`
	Wxids        string `json:"wxids"`
	Msg          string `json:"msg"`
	AutoNickname int    `json:"auto_nickname"`
}

type MsgSendAtResult struct{}

// RequestMsgSendAt 发送群艾特
func (a *Api) RequestMsgSendAt(data *MsgSendAtData) (*MsgSendAtResult, error) {
	res := &MsgSendAtResult{}
	if err := a.sendApi(MSG_SEND_AT, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type MsgSendCardData struct {
	Receiver   string `json:"receiver"`
	SharedWxid string `json:"shared_wxid"`
	Nickname   string `json:"nickname"`
}

type MsgSendCardResult struct{}

// 分享好友名片
func (a *Api) RequestMsgSendCard(data *MsgSendCardData) (*MsgSendCardResult, error) {
	res := &MsgSendCardResult{}
	if err := a.sendApi(MSG_SEND_CARD, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type MsgSendImageData struct {
	Receiver string `json:"receiver"`
	ImgPath  string `json:"img_path"`
}

type MsgSendImageResult struct{}

// 发送图片
func (a *Api) RequestMsgSendImage(data *MsgSendImageData) (*MsgSendImageResult, error) {
	res := &MsgSendImageResult{}
	if err := a.sendApi(MSG_SEND_IMAGE, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type MsgSendFileData struct {
	Receiver string `json:"receiver"`
	FilePath string `json:"file_path"`
}

type MsgSendFileResult struct{}

// 发送文件
func (a *Api) RequestMsgSendFile(data *MsgSendFileData) (*MsgSendFileResult, error) {
	res := &MsgSendFileResult{}
	if err := a.sendApi(MSG_SEND_FILE, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type MsgSendArticleData struct {
	Wxid     string `json:"wxid"`
	Title    string `json:"title"`
	Abstract string `json:"abstract"`
	Url      string `json:"url"`
	ImgPath  string `json:"img_path"`
}

type MsgSendArticleResult struct{}

// 发送xml文章
func (a *Api) RequestMsgSendArticle(data *MsgSendArticleData) (*MsgSendArticleResult, error) {
	res := &MsgSendArticleResult{}
	if err := a.sendApi(MSG_SEND_ARTICLE, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type MsgSendAppData struct {
	Wxid  string `json:"wxid"`
	Appid string `json:"appid"`
}

type MsgSendAppResult struct{}

// 发送小程序
func (a *Api) RequestMsgSendApp(data *MsgSendAppData) (*MsgSendAppResult, error) {
	res := &MsgSendAppResult{}
	if err := a.sendApi(MSG_SEND_APP, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type MsgStartHookData struct {
	Port int `json:"port"`
}

type CommonResult struct {
	Msg    int    `json:"msg"`
	Result string `json:"result"`
}

// 开启接收消息HOOK，只支持socket监听
func (a *Api) RequestMsgStartHook(data *MsgStartHookData) (*CommonResult, error) {
	res := &CommonResult{}
	if err := a.sendApi(MSG_START_HOOK, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type MsgStopHookResult struct {
	Msg    int    `json:"msg"`
	Result string `json:"result"`
}

// 关闭接收消息HOOK
func (a *Api) RequestMsgStopHook() (*CommonResult, error) {
	res := &CommonResult{}
	if err := a.sendApi(MSG_STOP_HOOK, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type MsgStartImageHookData struct {
	SavePath string `json:"save_path"`
}

// 开启图片消息HOOK
func (a *Api) RequestMsgStartImageHook(data *MsgStartImageHookData) (*CommonResult, error) {
	res := &CommonResult{}
	if err := a.sendApi(MSG_START_IMAGE_HOOK, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

// 关闭图片消息HOOK
func (a *Api) RequestMsgStopImageHook() (*CommonResult, error) {
	res := &CommonResult{}
	if err := a.sendApi(MSG_STOP_IMAGE_HOOK, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type MsgStartVoiceHookData struct {
	SavePath string `json:"save_path"`
}

// 开启语音消息HOOK
func (a *Api) RequestMsgStartVoiceHook(data *MsgStartVoiceHookData) (*CommonResult, error) {
	res := &CommonResult{}
	if err := a.sendApi(MSG_START_VOICE_HOOK, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

// 关闭语音消息HOOK
func (a *Api) RequestMsgStopVoiceHook() (*CommonResult, error) {
	res := &CommonResult{}
	if err := a.sendApi(MSG_STOP_VOICE_HOOK, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ContactGetListResult struct {
	Data   []ContractInfo `json:"data"`
	Result string         `json:"result"`
}

type ContractInfo struct {
	WxNickName   string `json:"wxNickName"`
	WxNumber     string `json:"wxNumber"`
	WxRemark     string `json:"wxRemark"`
	WxType       int    `json:"wxType"`
	WxVerifyFlag int    `json:"wxVerifyFlag"`
	Wxid         string `json:"wxid"`
}

// contact 获取联系人列表
func (a *Api) RequestContactGetList() (*ContactGetListResult, error) {
	res := &ContactGetListResult{}
	if err := a.sendApi(CONTACT_GET_LIST, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ContactCheckStatusData struct {
	Wxid string `json:"wxid"`
}

type ContactCheckStatusResult struct {
	Status int    `json:"status"`
	Result string `json:"result"`
}

// 检查是否被好友删除
func (a *Api) RequestContactCheckStatus(data *ContactCheckStatusData) (*ContactCheckStatusResult, error) {
	res := &ContactCheckStatusResult{}
	if err := a.sendApi(CONTACT_CHECK_STATUS, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ContactDelData struct {
	Wxid string `json:"wxid"`
}

type ContactDelResult struct {
	Msg    int    `json:"msg"`
	Result string `json:"result"`
}

// 删除好友
func (a *Api) RequestContactDel(data *ContactDelData) (*ContactDelResult, error) {
	res := &ContactDelResult{}
	if err := a.sendApi(CONTACT_DEL, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ContactSearchByCacheData struct {
	Wxid string `json:"wxid"`
}

type ContactSearchByCacheResult struct{}

// 从内存中获取好友信息
func (a *Api) RequestContactSearchByCache(data *ContactSearchByCacheData) (*ContactSearchByCacheResult, error) {
	res := &ContactSearchByCacheResult{}
	if err := a.sendApi(CONTACT_SEARCH_BY_CACHE, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ContactSearchByNetData struct {
	Keyword string `json:"keyword"`
}

type ContactSearchByNetResult struct{}

// 网络搜索用户信息
func (a *Api) RequestContactSearchByNet(data *ContactSearchByNetData) (*ContactSearchByNetResult, error) {
	res := &ContactSearchByNetResult{}
	if err := a.sendApi(CONTACT_SEARCH_BY_NET, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ContactAddByWxidData struct {
	Wxid string `json:"wxid"`
	Msg  string `json:"msg"`
}

type ContactAddByWxidResult struct{}

// wxid加好友
func (a *Api) RequestContactAddByWxid(data *ContactAddByWxidData) (*ContactAddByWxidResult, error) {
	res := &ContactAddByWxidResult{}
	if err := a.sendApi(CONTACT_ADD_BY_WXID, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ContactAddByV3Data struct {
	V3      string `json:"v3"`
	Msg     string `json:"msg"`
	AddType int    `json:"add_type"`
}

type ContactAddByV3Result struct{}

// v3数据加好友
func (a *Api) RequestContactAddByV3(data *ContactAddByV3Data) (*ContactAddByV3Result, error) {
	res := &ContactAddByV3Result{}
	if err := a.sendApi(CONTACT_ADD_BY_V3, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ContactAddByPublicIdData struct {
	PublicId string `json:"public_id"`
}

type ContactAddByPublicIdResult struct{}

// 关注公众号
func (a *Api) RequestContactAddByPublicId(data *ContactAddByPublicIdData) (*ContactAddByPublicIdResult, error) {
	res := &ContactAddByPublicIdResult{}
	if err := a.sendApi(CONTACT_ADD_BY_PUBLIC_ID, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ContactVerifyApplyData struct {
	V3 string `json:"v3"`
	V4 string `json:"v4"`
}

type ContactVerifyApplyResult struct{}

// 通过好友请求
func (a *Api) RequestContactVerifyApply(data *ContactVerifyApplyData) (*ContactVerifyApplyResult, error) {
	res := &ContactVerifyApplyResult{}
	if err := a.sendApi(CONTACT_VERIFY_APPLY, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ContactEditRemarkData struct {
	Wxid   string `json:"wxid"`
	Remark string `json:"remark"`
}

type ContactEditRemarkResult struct{}

// 修改备注
func (a *Api) RequestContactEditRemark(data *ContactEditRemarkData) (*ContactEditRemarkResult, error) {
	res := &ContactEditRemarkResult{}
	if err := a.sendApi(CONTACT_EDIT_REMARK, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ChatroomGetMemberListData struct {
	ChatroomId string `json:"chatroom_id"`
}

type ChatroomGetMemberListResult struct {
	Members string `json:"members"`
	Result  string `json:"result"`
}

// 获取群成员列表
func (a *Api) RequestChatroomGetMemberList(data *ChatroomGetMemberListData) (*ChatroomGetMemberListResult, error) {
	res := &ChatroomGetMemberListResult{}
	if err := a.sendApi(CHATROOM_GET_MEMBER_LIST, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ChatroomGetMemberNicknameData struct {
	ChatroomId string `json:"chatroom_id"`
	Wxid       string `json:"wxid"`
}

type ChatroomGetMemberNicknameResult struct {
	Nickname string `json:"nickname"`
	Result   string `json:"result"`
}

// 获取指定群成员昵称
func (a *Api) RequestChatroomGetMemberNickname(data *ChatroomGetMemberNicknameData) (*ChatroomGetMemberNicknameResult, error) {
	res := &ChatroomGetMemberNicknameResult{}
	if err := a.sendApi(CHATROOM_GET_MEMBER_NICKNAME, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ChatroomDelMemberData struct {
	ChatroomId string `json:"chatroom_id"`
	Wxids      string `json:"wxids"`
}

type ChatroomDelMemberResult struct {
	Msg    int    `json:"msg"`
	Result string `json:"result"`
}

// 删除群成员
func (a *Api) RequestChatroomDelMember(data *ChatroomDelMemberData) (*ChatroomDelMemberResult, error) {
	res := &ChatroomDelMemberResult{}
	if err := a.sendApi(CHATROOM_DEL_MEMBER, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ChatroomAddMemberData struct {
	ChatroomId string `json:"chatroom_id"`
	Wxids      string `json:"wxids"`
}

type ChatroomAddMemberResult struct {
	Msg    int    `json:"msg"`
	Result string `json:"result"`
}

// 添加群成员
func (a *Api) RequestChatroomAddMember(data *ChatroomAddMemberData) (*ChatroomAddMemberResult, error) {
	res := &ChatroomAddMemberResult{}
	if err := a.sendApi(CHATROOM_ADD_MEMBER, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ChatroomSetAnnouncementData struct {
	ChatroomId   string `json:"chatroom_id"`
	Announcement string `json:"announcement"`
}

type ChatroomSetAnnouncementResult struct {
	Msg    int    `json:"msg"`
	Result string `json:"result"`
}

// 设置群公告
func (a *Api) RequestChatroomSetAnnouncement(data *ChatroomSetAnnouncementData) (*ChatroomSetAnnouncementResult, error) {
	res := &ChatroomSetAnnouncementResult{}
	if err := a.sendApi(CHATROOM_SET_ANNOUNCEMENT, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ChatroomSetChatroomNameData struct {
	ChatroomId   string `json:"chatroom_id"`
	ChatroomName string `json:"chatroom_name"`
}

type ChatroomSetChatroomNameResult struct {
	Msg    int    `json:"msg"`
	Result string `json:"result"`
}

// 设置群聊名称
func (a *Api) RequestChatroomSetChatroomName(data *ChatroomSetChatroomNameData) (*ChatroomSetChatroomNameResult, error) {
	res := &ChatroomSetChatroomNameResult{}
	if err := a.sendApi(CHATROOM_SET_CHATROOM_NAME, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ChatroomSetSelfNicknameData struct {
	ChatroomId string `json:"chatroom_id"`
	Nickname   string `json:"nickname"`
}

type ChatroomSetSelfNicknameResult struct{}

// 设置群内个人昵称
func (a *Api) RequestChatroomSetSelfNickname(data *ChatroomSetSelfNicknameData) (*ChatroomSetSelfNicknameResult, error) {
	res := &ChatroomSetSelfNicknameResult{}
	if err := a.sendApi(CHATROOM_SET_SELF_NICKNAME, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type DatabaseGetHandlesResult struct{}

// 获取数据库句柄
func (a *Api) RequestDatabaseGetHandles() (*DatabaseGetHandlesResult, error) {
	res := &DatabaseGetHandlesResult{}
	if err := a.sendApi(DATABASE_GET_HANDLES, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type DatabaseBackupData struct {
	DbHandle int    `json:"db_handle"`
	SavePath string `json:"save_path"`
}

type DatabaseBackupResult struct{}

// 备份数据库
func (a *Api) RequestDatabaseBackup(data *DatabaseBackupData) (*DatabaseBackupResult, error) {
	res := &DatabaseBackupResult{}
	if err := a.sendApi(DATABASE_BACKUP, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type DatabaseQueryData struct {
	DbHandle int    `json:"db_handle"`
	Sql      string `json:"sql"`
}

type DatabaseQueryResult struct{}

// 数据库查询
func (a *Api) RequestDatabaseQuery(data *DatabaseQueryData) (*DatabaseQueryResult, error) {
	res := &DatabaseQueryResult{}
	if err := a.sendApi(DATABASE_QUERY, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type SetVersionData struct {
	Version string `json:"version"`
}

type SetVersionResult struct{}

// 修改微信版本号
func (a *Api) RequestSetVersion(data *SetVersionData) (*SetVersionResult, error) {
	res := &SetVersionResult{}
	if err := a.sendApi(SET_VERSION, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type LogStartHookResult struct{}

// 开启日志信息HOOK
func (a *Api) RequestLogStartHook() (*LogStartHookResult, error) {
	res := &LogStartHookResult{}
	if err := a.sendApi(LOG_START_HOOK, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type LogStopHookResult struct{}

// 关闭日志信息HOOK
func (a *Api) RequestLogStopHook() (*LogStopHookResult, error) {
	res := &LogStopHookResult{}
	if err := a.sendApi(LOG_STOP_HOOK, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type BrowserOpenWithUrlData struct {
	Url string `json:"url"`
}

type BrowserOpenWithUrlResult struct{}

// 打开微信内置浏览器
func (a *Api) RequestBrowserOpenWithUrl(data *BrowserOpenWithUrlData) (*BrowserOpenWithUrlResult, error) {
	res := &BrowserOpenWithUrlResult{}
	if err := a.sendApi(BROWSER_OPEN_WITH_URL, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type GetPublicMsgData struct {
	PublicId string `json:"public_id"`
	Offset   string `json:"offset"`
}

type GetPublicMsgResult struct{}

// 获取公众号历史消息
func (a *Api) RequestGetPublicMsg(data *GetPublicMsgData) (*GetPublicMsgResult, error) {
	res := &GetPublicMsgResult{}
	if err := a.sendApi(GET_PUBLIC_MSG, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type MsgForwardMessageData struct {
	Wxid  string `json:"wxid"`
	Msgid int    `json:"msgid"`
}

type MsgForwardMessageResult struct{}

// 转发消息
func (a *Api) RequestMsgForwardMessage(data *MsgForwardMessageData) (*MsgForwardMessageResult, error) {
	res := &MsgForwardMessageResult{}
	if err := a.sendApi(MSG_FORWARD_MESSAGE, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type GetQrcodeImageResult struct{}

// 获取二维码
func (a *Api) RequestGetQrcodeImage() (*GetQrcodeImageResult, error) {
	res := &GetQrcodeImageResult{}
	if err := a.sendApi(GET_QRCODE_IMAGE, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type GetA8keyData struct {
	Url string `json:"url"`
}

type GetA8keyResult struct{}

// 获取A8Key
func (a *Api) RequestGetA8key(data *GetA8keyData) (*GetA8keyResult, error) {
	res := &GetA8keyResult{}
	if err := a.sendApi(GET_A8KEY, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type MsgSendXmlData struct {
	Wxid    string `json:"wxid"`
	Xml     string `json:"xml"`
	ImgPath string `json:"img_path"`
}

type MsgSendXmlResult struct{}

// 发送xml消息
func (a *Api) RequestMsgSendXml(data *MsgSendXmlData) (*MsgSendXmlResult, error) {
	res := &MsgSendXmlResult{}
	if err := a.sendApi(MSG_SEND_XML, data, res); err != nil {
		return nil, err
	}

	return res, nil
}

type LogoutResult struct{}

// 退出登录
func (a *Api) RequestLogout() (*LogoutResult, error) {
	res := &LogoutResult{}
	if err := a.sendApi(LOGOUT, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (a *Api) sendApi(t Type, data, v interface{}) error {
	js := []byte(`{}`)
	if data != nil {
		js, _ = json.Marshal(data)
	}
	host := fmt.Sprintf(a.host, t)
	cli := http.DefaultClient
	cli.Timeout = a.timeout
	req, err := http.NewRequest("POST", host, bytes.NewReader(js))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(respData, v); err != nil {
		return err
	}
	return nil
}

type TextMessage struct {
	Extrainfo string `json:"extrainfo"`
	Filepath  string `json:"filepath"`
	IsSendMsg int    `json:"isSendMsg"`
	Message   string `json:"message"`
	Msgid     int64  `json:"msgid"`
	Pid       int    `json:"pid"`
	Self      string `json:"self"`
	Sender    string `json:"sender"`
	Sign      string `json:"sign"`
	ThumbPath string `json:"thumb_path"`
	Time      string `json:"time"`
	Timestamp int    `json:"timestamp"`
	Type      int    `json:"type"`
	Wxid      string `json:"wxid"`
}
