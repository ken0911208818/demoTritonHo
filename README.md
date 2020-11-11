# demoTritonHo
練習建構一個基本的CRUD fork [TritonHo - demo ](https://github.com/TritonHo/demo)

## 啟動
```bash=
 /bin/zsh demoTritonHo/dev_env.sh
```
## 基本http要素
- 路由
- 驗證
- 商業邏輯
- 資料庫


## 使用套件

- mux
- gorm
- jwt-go
- pgsql
- uuid
- xorm


## 資料庫差異

使用`xorm` 比較接近原生ORM`sql`語法 在使用`gorm`時查詢了許多`orm`的寫法 為了比較原作者`(xorm)`與`gorm`使用上的差異花了滿多時間比較實作差異 

> 某些場合使用原生語法比用ORM還方便許多

> 再新增使用者判斷是否存在
> 如果使用ORM 要先判斷是否select有值 再進行新增 若使用原生語法就能query一次結束
```sql=
insert into users(id, email, password_digest, first_name, last_name)
select ?, ?, ?, ?, ?
where not exists (select 1 from users where email = ?)
```

> maybe can use gorm.firstorCreate() function to check
## 中介層

原作者習慣laravel方便的middleware 寫go時非常不順手，一開始容易卡住為何要這麼做，以及在寫router 時都要呼叫 `http.RequestWiter` `http.Request` 等等`指數`的概念 後來多方嘗試後 才知道傳入 `*http.Request` 是為了讓middleware判斷完應有的行為後進行 return http.code or go ahead 

## 驗證 

還沒有看懂 待捕

## mux 

> 一種路由器 可以方便存取 http 參數 進行商業邏輯的判斷

```golang
router := mux.NewRouter() // 宣告使用mux router 
router.HandleFunc("/v1/users/{userId:"+uuidRegexp+"}", middleware.Wrap(handler.UserUpdate)).Methods("PUT") // 可以存取userId 的值 也可以在上面加上簡易的驗證條件做第一部篩選
mux.Vars(*http.Request) // 可以將http宣告的變數存取
```

## jwt 

> 使用電腦上的私鑰 並以RS512 的方式加密 

## 觀念技巧

> 能將create 或是 update 這種對於資料庫的操作寫成function 可以減少 if err!=nil 的次數 在內部回傳 db create 所需要的參數即可