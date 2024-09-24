package router

import (
	"myhomesv/internal/controller"
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

////////////////////////////////////////////////////////////////////////////
// router.go - 家計簿サービスのルート設定
// このファイルには、家計簿サービスの各画面に対応するルートとハンドラ関数が含まれています。
////////////////////////////////////////////////////////////////////////////

// 認証ミドルウェア
func authMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}
	c.Next()
}

// Gin はデフォルトでは、1つの html.Template しか使用できません。 go 1.6 の block template のような機能が使用できる a multitemplate render を検討してください。
// https://gin-gonic.com/ja/docs/examples/multiple-template/
func createRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	// login
	r.AddFromFiles("login", "web/templates/base.html", "web/templates/user/login.html")
	r.AddFromFiles("request-password-reset", "web/templates/base.html", "web/templates/user/request-password-reset.html")
	r.AddFromFiles("reset-password", "web/templates/base.html", "web/templates/user/reset-password.html")
	// dashboard
	r.AddFromFiles("dashboard", "web/templates/base.html", "web/templates/dashboard.html")
	// households
	r.AddFromFiles("wallets", "web/templates/base.html", "web/templates/households/wallets.html")
	r.AddFromFiles("categories", "web/templates/base.html", "web/templates/households/categories.html")
	r.AddFromFiles("reports", "web/templates/base.html", "web/templates/households/reports.html")
	r.AddFromFiles("budget", "web/templates/base.html", "web/templates/households/budget.html")
	r.AddFromFiles("settings", "web/templates/base.html", "web/templates/households/settings.html")

	return r
}

// SetupRouter - ルートを設定し、Ginエンジンを返します。
func SetupRouter(bc controller.IBudgetController) *gin.Engine {
	r := gin.Default()

	// セッションストアを設定
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge:   3600 * 3, // セッションの有効期限を秒単位で設定（例: 3600秒 = 1時間）
		HttpOnly: true,     // クライアントサイドのJavaScriptからセッションをアクセスできないようにする
		Secure:   true,     // HTTPSを使用している場合にのみセッションを送信する
	})
	r.Use(sessions.Sessions("mysession", store))

	// TODO: CSRFミドルウェアを設定
	// r.Use(csrf.Middleware(csrf.Options{
	// 	Secret: "secret",
	// 	ErrorFunc: func(c *gin.Context) {
	// 		c.String(400, "CSRF token mismatch")
	// 		c.Abort()
	// 	},
	// }))
	// TODO: CSPヘッダーを設定
	// r.Use(func(c *gin.Context) {
	// 	c.Header("Content-Security-Policy", "default-src 'self'")
	// 	c.Header("X-Content-Type-Options", "nosniff")
	// 	c.Header("X-Frame-Options", "DENY")
	// 	c.Header("X-XSS-Protection", "1; mode=block")
	// 	c.Next()
	// })
	// r.LoadHTMLGlob("web/templates/*")
	r.HTMLRender = createRender()
	// 画面遷移 ===========================================
	// ルート: ホーム画面
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})
	// User: ログイン画面
	r.GET("/login", func(c *gin.Context) {
		//TODO: CSRFトークンを取得
		//c.HTML(200, "login", gin.H{"title": "login", "csrfToken": csrf.GetToken(c)})
		c.HTML(200, "login", gin.H{"title": "login"})
	})
	// User: リセットパスワード要求画面
	r.GET("/request-password-reset", func(c *gin.Context) {
		c.HTML(200, "request-password-reset", gin.H{"title": "request-password-reset"})
	})
	// User: リセットパスワード画面
	r.GET("/reset-password", func(c *gin.Context) {
		c.HTML(200, "reset-password", gin.H{"title": "reset-password", "token": c.Query("token")})
	})
	// User: ダッシュボード画面
	r.GET("/dashboard", authMiddleware, func(c *gin.Context) {
		c.HTML(200, "dashboard", gin.H{"title": "dashboard"})
	})
	// User: 複数財布管理画面
	r.GET("/wallets", authMiddleware, func(c *gin.Context) {
		c.HTML(200, "wallets", gin.H{"title": "wallets"})
	})
	// ルート: 収支データインポート画面
	r.GET("/import", authMiddleware, func(c *gin.Context) {
		c.HTML(200, "wallets", gin.H{"title": "wallets"})
	})
	// ルート: 費目ベース収支管理画面
	r.GET("/categories", authMiddleware, func(c *gin.Context) {
		c.HTML(200, "categories", gin.H{"title": "categories"})
	})
	// ルート: レポート画面
	r.GET("/reports", authMiddleware, func(c *gin.Context) {
		c.HTML(200, "reports", gin.H{"title": "reports"})
	})
	// ルート: 予算管理画面
	r.GET("/budget", authMiddleware, func(c *gin.Context) {
		c.HTML(200, "budget", gin.H{"title": "budget"})
	})
	// ルート: 設定画面
	r.GET("/settings", authMiddleware, func(c *gin.Context) {
		c.HTML(200, "settings", gin.H{"title": "settings"})
	})
	// ログイン処理 ===========================================
	// ログイン処理: ログイン
	r.POST("/login", func(c *gin.Context) {
		bc.Login(c)
		if c.Writer.Status() == http.StatusOK {
			session := sessions.Default(c)
			session.Clear()                 // 既存のセッションデータをクリア
			session.Set("user", "username") // ここでユーザー情報をセッションに保存
			if err := session.Save(); err != nil {
				c.HTML(500, "login", gin.H{"title": "login", "LoginError": "セッションの保存に失敗しました"})
				return
			}
			c.Redirect(http.StatusFound, "/dashboard")
		} else {
			c.HTML(200, "login", gin.H{"title": "login", "LoginError": "ログインに失敗しました"})
		}
	})
	// ログイン処理: サインイン
	r.POST("/signin", func(c *gin.Context) {
		bc.SignUp(c)
		if c.Writer.Status() == http.StatusCreated {
			c.HTML(200, "login", gin.H{"title": "login", "SignUpSuccess": "新規ユーザー登録に成功しました ログインして下さい"})
		} else {
			c.HTML(500, "login", gin.H{"title": "login", "SignUpError": "新規ユーザー登録に失敗しました"})
		}
	})
	// ログイン処理: パスワードリセット要求
	r.POST("/request-password-reset", func(c *gin.Context) {
		bc.RequestSentResetEmail(c)
		if c.Writer.Status() == http.StatusOK {
			c.HTML(200, "request-password-reset", gin.H{"title": "request-password-reset", "RequestPasswordResetSuccess": "パスワードリセット要求に成功しました メールを確認して下さい"})
		} else {
			c.HTML(500, "request-password-reset", gin.H{"title": "request-password-reset", "RequestPasswordResetError": "パスワードリセット要求に失敗しました"})
		}
	})
	// ログイン処理: パスワードリセット
	r.POST("/reset-password", func(c *gin.Context) {
		bc.ResetPassword(c)
		if c.Writer.Status() == http.StatusOK {
			c.HTML(200, "reset-password", gin.H{"title": "reset-password", "PasswordResetSuccess": "パスワードリセットに成功しました ログインして下さい"})
		} else {
			c.HTML(500, "reset-password", gin.H{"title": "reset-password", "PasswordResetError": "パスワードリセットに失敗しました"})
		}
	})
	return r
}
