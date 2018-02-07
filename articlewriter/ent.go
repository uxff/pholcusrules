package model

import "time"

type ArticleEntity struct {
	Id            int       `json:"id"          xorm:"pk autoincr"`
	Title         string    `json:"title"       xorm:"not null"`
	Type          int       `json:"type"        xorm:"not null"`
	Content       string    `json:"content"     xorm:"not null"`
	Outer_url     string    `json:"outer_url"   xorm:"not null default ''"`
	Visit_url     string    `json:"visit_url"   xorm:"not null default ''"`
	Surface_url   string    `json:"surface_url" xorm:"not null default ''"`
	Abstract      string    `json:"abstract"    xorm:"not null default ''"`
	Status        int       `json:"status"      xorm:"not null default '0'"`
	Remark        string    `json:"remark"      xorm:"not null default ''"`
	View_count    int       `json:"view_count"  xorm:"not null default '0'"`
	Share_count   int       `json:"share_count" xorm:"not null"`
	Favor_count   int       `json:"favor_count" xorm:"not null"`
	Comment_count int       `json:"comment_count" xorm:"not null"`
	Pubdate       time.Time `json:"pubdate"     xorm:"not null default '0000-00-00 00:00:00'"`
	Create_time   time.Time `json:"create_time"   xorm:"not null default '0000-00-00 00:00:00'"`
	Last_modified time.Time `json:"last_modified" xorm:"not null default '0000-00-00 00:00:00'"`
	Admin_id      int       `json:"admin_id"   xorm:"not null default '1'"`
	Admin_name    string    `json:"admin_name" xorm:"not null default ''"`
	Origin        string    `json:"origin"     xorm:"not null default ''"`
	Author        string    `json:"author"     xorm:"not null default ''"`
}

func (this ArticleEntity) TableName() string {
	return "fh_article"
}
