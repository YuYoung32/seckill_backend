package handler

import (
	. "context"
	"errors"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"net/http"
	"sk_user_srv/conf"
	"sk_user_srv/database"
	. "sk_user_srv/proto/user"
	. "sk_user_srv/utils"

	"strconv"
	"time"
)

type UserImpl struct {
}

var codeCache *cache.Cache

func init() {
	config := conf.GetCacheConf()
	codeCache = cache.New(time.Duration(config.ExpireTime)*time.Second, time.Duration(config.CleanTime)*time.Second)
}

func (u *UserImpl) Register(ctx Context, in *RegisterUserRequest, out *GeneralResponse) error {
	email := in.User.BasicInfo.Email
	trueCode, ok := codeCache.Get(email)

	//验证码是否匹配
	if !ok || in.Code != trueCode {
		logrus.WithField("email", email).Info("错误的或失效的验证码")
		out.Msg = "错误的或失效验证码"
		out.Code = strconv.Itoa(http.StatusNotFound)
		return errors.New("错误的或失效的验证码")
	}

	db := database.GetDBConn()
	if res := db.Where("username=?", email).First(&database.User{}); res.RowsAffected > 0 {
		logrus.WithField("email", email).Info("用户已存在")
		out.Msg = "用户已存在"
		out.Code = strconv.Itoa(http.StatusNotFound)
		return nil
	}

	hashedPass, err := Encryption(in.User.BasicInfo.Password)
	if err != nil {
		logrus.WithField("email", email).Info("密码加密失败")
		out.Msg = "密码加密失败"
		out.Code = strconv.Itoa(http.StatusInternalServerError)
		return nil
	}

	db.Create(&database.User{
		Username:    in.User.Username,
		Email:       email,
		Password:    hashedPass,
		Description: in.User.Description,
		Status:      in.User.Status,
	})

	logrus.WithField("email", email).Info("注册成功")
	out.Msg = "恭喜您注册成功"
	out.Code = strconv.Itoa(http.StatusOK)
	return nil
}

func (u *UserImpl) Login(ctx Context, in *GeneralRequest, out *GeneralResponse) error {
	email := in.User.Email
	password := in.User.Password
	db := database.GetDBConn()

	user := database.User{}
	if res := db.Where("username=?", email).First(&user); res.RowsAffected == 0 {
		logrus.WithField("email", email).Info("用户不存在")
		out.Msg = "用户不存在"
		out.Code = strconv.Itoa(http.StatusNotFound)
		return errors.New("用户不存在")
	}

	if ValidationPassword(user.Password, password) {
		logrus.WithField("email", email).Info("登录成功")
		out.Msg = "登录成功"
		out.Code = strconv.Itoa(http.StatusOK)
	} else {
		logrus.WithField("email", email).WithField("password", password).Info("密码错误")
		out.Msg = "密码错误"
		out.Code = strconv.Itoa(http.StatusNotFound)
		return errors.New("密码错误")
	}
	return nil
}

func (u *UserImpl) AdminLogin(ctx Context, in *GeneralRequest, out *GeneralResponse) error {
	email := in.User.Email
	password := in.User.Password
	db := database.GetDBConn()

	//管理员登陆
	if len(email) >= 7 && email[:7] == "skadmin" {
		admin := database.Admin{}
		db.Where("username=?", email).First(&admin)
		if len(admin.Password) > 0 && ValidationPassword(admin.Password, password) {
			logrus.WithField("username", email).Info("管理员登录成功")
			out.Msg = "管理员登录成功"
			out.Code = strconv.Itoa(http.StatusOK)
			return nil
		}
	}
	logrus.WithField("username", email).WithField("password", password).Info("管理员密码错误")
	out.Msg = "密码错误"
	out.Code = strconv.Itoa(http.StatusNotFound)
	return errors.New("密码错误")
}

func (u *UserImpl) SendEmail(ctx Context, in *GeneralRequest, out *GeneralResponse) error {
	randCode := GenEmailCode(4)
	email := in.User.Email
	config := conf.GetEmailConf()

	codeCache.Set(email, randCode, cache.DefaultExpiration)

	content := "秒杀平台：您请求的验证码是" + randCode + "，有效期3分钟。"
	err := SendEmail(email, content, config.Username, config.Password, config.Host, config.Port)
	if err != nil {
		logrus.WithField("module", "send_email").Error(err)
		out.Msg = "发送邮件失败"
		out.Code = strconv.Itoa(http.StatusInternalServerError)
		return err
	}

	logrus.WithFields(logrus.Fields{
		"sender": config.Username,
		"email":  email,
		"code":   randCode,
	}).Info("发送邮件成功")

	out.Code = strconv.Itoa(http.StatusOK)
	out.Msg = "发送邮件成功，请查收！"
	return nil
}

func (u *UserImpl) GetUserInfo(ctx Context, in *GetUserInfoRequest, out *GetUserInfoResponse) error {
	db := database.GetDBConn()
	start := in.Start
	amount := in.Amount

	var res []database.User
	var total int
	e := db.Limit(amount).Offset(start).Find(&res)
	e = db.Model(&database.User{}).Count(&total)
	if e.Error != nil {
		out.General = &GeneralResponse{
			Code: strconv.Itoa(http.StatusInternalServerError),
			Msg:  "数据库查询错误",
		}
		return e.Error
	}

	var userInfo []*UserInfo
	for _, r := range res {
		userInfo = append(userInfo, &UserInfo{
			BasicInfo: &BasicUserInfo{
				Email:    r.Email,
				Password: "", //此处不应传密码
			},
			Username:    r.Username,
			Description: r.Description,
			Status:      r.Status,
			CreateTime:  r.CreatedAt.String(),
		})
	}

	//这样不行，不能替换指针，只能在原来的基础上进行添加东西
	*out = GetUserInfoResponse{
		General: &GeneralResponse{
			Code: strconv.Itoa(http.StatusOK),
			Msg:  "查询成功",
		},
		User: userInfo,
	}
	//当然也不能直接给General赋值，因为这个指针是没有分配空间的
	//out.General.Code="200"

	out.General = &GeneralResponse{
		Code: strconv.Itoa(http.StatusOK),
		Msg:  "查询成功",
	}
	out.User = userInfo
	out.Total = int32(total)

	return nil
}
