package bilive

type Message struct {
	CMD CMD `json:"cmd"`
}

type CMD string

const (

	// 弹幕类

	CMD_DANMU_MSG              CMD = "DANMU_MSG"              // 弹幕消息
	CMD_WELCOME_GUARD          CMD = "WELCOME_GUARD"          //欢迎xxx老爷
	CMD_ENTRY_EFFECT           CMD = "ENTRY_EFFECT"           // 欢迎舰长进入房间
	CMD_WELCOME                CMD = "WELCOME"                // 欢迎xxx进入房间
	CMD_SUPER_CHAT_MESSAGE_JPN CMD = "SUPER_CHAT_MESSAGE_JPN" // 二个都是SC留言
	CMD_SUPER_CHAT_MESSAGE     CMD = "SUPER_CHAT_MESSAGE"     // 二个都是SC留言

	// 礼物类

	CMD_SEND_GIFT  CMD = "SEND_GIFT"  // 投喂礼物
	CMD_COMBO_SEND CMD = "COMBO_SEND" // 连击礼物

	// 天选之人类

	CMD_ANCHOR_LOT_START CMD = "ANCHOR_LOT_START" // 天选之人开始完整信息
	CMD_ANCHOR_LOT_END   CMD = "ANCHOR_LOT_END"   // 天选之人获奖id
	CMD_ANCHOR_LOT_AWARD CMD = "ANCHOR_LOT_AWARD" // 天选之人获奖完整信息

	// 上船类

	CMD_GUARD_BUY      CMD = "GUARD_BUY"      // 上舰长
	CMD_USER_TOAST_MSG CMD = "USER_TOAST_MSG" // 续费了舰长
	CMD_NOTICE_MSG     CMD = "NOTICE_MSG"     // 在本房间续费了舰长

	// 分区排行类

	CMD_ACTIVITY_BANNER_UPDATE_V2 CMD = "ACTIVITY_BANNER_UPDATE_V2" // 小时榜变动

	// 关注数变化类

	CMD_ROOM_REAL_TIME_MESSAGE_UPDATE CMD = "ROOM_REAL_TIME_MESSAGE_UPDATE" // 上舰长

)
