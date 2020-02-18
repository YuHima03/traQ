package v3

import (
	"github.com/labstack/echo/v4"
	"github.com/leandro-lugaresi/hub"
	"github.com/traPtitech/traQ/rbac"
	"github.com/traPtitech/traQ/rbac/permission"
	"github.com/traPtitech/traQ/realtime"
	"github.com/traPtitech/traQ/realtime/ws"
	"github.com/traPtitech/traQ/repository"
	"github.com/traPtitech/traQ/router/middlewares"
	"go.uber.org/zap"
)

type Handlers struct {
	RBAC     rbac.RBAC
	Repo     repository.Repository
	WS       *ws.Streamer
	Hub      *hub.Hub
	Logger   *zap.Logger
	Realtime *realtime.Service

	Version  string
	Revision string
}

// Setup APIルーティングを行います
func (h *Handlers) Setup(e *echo.Group) {
	// middleware preparation
	requires := middlewares.AccessControlMiddlewareGenerator(h.RBAC)
	retrieve := middlewares.NewParamRetriever(h.Repo)

	api := e.Group("/v3", middlewares.UserAuthenticate(h.Repo))
	{
		apiUsers := api.Group("/users")
		{
			apiUsers.GET("", NotImplemented)
			apiUsers.POST("", NotImplemented)
			apiUsersUID := apiUsers.Group("/:userID", retrieve.UserID(false))
			{
				apiUsersUID.GET("", NotImplemented)
				apiUsersUID.PATCH("", NotImplemented)
				apiUsersUID.POST("/messages", NotImplemented)
				apiUsersUID.GET("/messages", NotImplemented)
				apiUsersUID.GET("/icon", NotImplemented)
				apiUsersUID.PUT("/icon", NotImplemented)
				apiUsersUID.PUT("/password", NotImplemented)
				apiUsersUIDTags := apiUsersUID.Group("/tags")
				{
					apiUsersUIDTags.GET("", NotImplemented)
					apiUsersUIDTags.POST("", NotImplemented)
					apiUsersUIDTagsTID := apiUsersUIDTags.Group("/:tagID")
					{
						apiUsersUIDTagsTID.PATCH("", NotImplemented)
						apiUsersUIDTagsTID.DELETE("", NotImplemented)
					}
				}
			}
			apiUsersMe := apiUsers.Group("/me")
			{
				apiUsersMe.GET("", NotImplemented)
				apiUsersMe.PATCH("", NotImplemented)
				apiUsersMe.GET("/stamp-history", NotImplemented)
				apiUsersMe.GET("/qr-code", h.GetMyQRCode)
				apiUsersMe.GET("/subscription", NotImplemented)
				apiUsersMe.PUT("/subscription/:channelID", NotImplemented)
				apiUsersMe.GET("/icon", NotImplemented)
				apiUsersMe.PUT("/icon", NotImplemented)
				apiUsersMe.PUT("/password", h.PutMyPassword)
				apiUsersMe.POST("/fcm-device", NotImplemented)
				apiUsersMeTags := apiUsersMe.Group("/tags")
				{
					apiUsersMeTags.GET("", NotImplemented)
					apiUsersMeTags.POST("", NotImplemented)
					apiUsersMeTagsTID := apiUsersMeTags.Group("/:tagID")
					{
						apiUsersMeTagsTID.PATCH("", NotImplemented)
						apiUsersMeTagsTID.DELETE("", NotImplemented)
					}
				}
				apiUsersMeStars := apiUsersMe.Group("/stars")
				{
					apiUsersMeStars.GET("", NotImplemented)
					apiUsersMeStars.POST("", NotImplemented)
					apiUsersMeStars.DELETE("/:channelID", NotImplemented)
				}
				apiUsersMe.GET("/unread", NotImplemented)
				apiUsersMe.DELETE("/unread", NotImplemented)
				apiUsersMe.GET("/sessions", NotImplemented)
				apiUsersMe.DELETE("/sessions/:sessionID", NotImplemented)
				apiUsersMe.GET("/tokens", NotImplemented)
				apiUsersMe.DELETE("/tokens/:tokenID", NotImplemented)
			}
		}
		apiChannels := api.Group("/channels")
		{
			apiChannels.GET("", NotImplemented)
			apiChannels.POST("", NotImplemented)
			apiChannelsCID := apiChannels.Group("/:channelID")
			{
				apiChannelsCID.GET("", NotImplemented)
				apiChannelsCID.PATCH("", NotImplemented)
				apiChannelsCID.GET("/messages", NotImplemented)
				apiChannelsCID.POST("/messages", NotImplemented)
				apiChannelsCID.GET("/stats", NotImplemented)
				apiChannelsCID.GET("/topic", NotImplemented)
				apiChannelsCID.PUT("/topic", NotImplemented)
				apiChannelsCID.GET("/viewers", NotImplemented)
				apiChannelsCID.GET("/pins", NotImplemented)
				apiChannelsCID.GET("/subscribers", NotImplemented)
				apiChannelsCID.PUT("/subscribers", NotImplemented)
				apiChannelsCID.PATCH("/subscribers", NotImplemented)
				apiChannelsCID.GET("/bots", NotImplemented)
				apiChannelsCID.GET("/events", NotImplemented)
			}
		}
		apiMessages := api.Group("/messages")
		{
			apiMessagesMID := apiMessages.Group("/:messageID")
			{
				apiMessagesMID.GET("", NotImplemented)
				apiMessagesMID.PUT("", NotImplemented)
				apiMessagesMID.DELETE("", NotImplemented)
				apiMessagesMID.GET("/pin", NotImplemented)
				apiMessagesMID.POST("/pin", NotImplemented)
				apiMessagesMID.DELETE("/pin", NotImplemented)
				apiMessagesMIDStamps := apiMessagesMID.Group("/stamps")
				{
					apiMessagesMIDStamps.GET("", NotImplemented)
					apiMessagesMIDStampsSID := apiMessagesMIDStamps.Group("/:stampID")
					{
						apiMessagesMIDStampsSID.POST("", NotImplemented)
						apiMessagesMIDStampsSID.DELETE("", NotImplemented)
					}
				}
			}
		}
		apiFiles := api.Group("/files")
		{
			apiFiles.GET("", NotImplemented)
			apiFiles.POST("", NotImplemented)
			apiFilesFID := apiFiles.Group("/:fileID")
			{
				apiFilesFID.GET("", NotImplemented)
				apiFilesFID.DELETE("", NotImplemented)
				apiFilesFID.GET("/meta", NotImplemented)
				apiFilesFID.GET("/thumbnail", NotImplemented)
			}
		}
		apiTags := api.Group("/tags")
		{
			apiTagsTID := apiTags.Group("/:tagID")
			{
				apiTagsTID.GET("", NotImplemented)
			}
		}
		apiStamps := api.Group("/stamps")
		{
			apiStamps.GET("", NotImplemented)
			apiStamps.POST("", NotImplemented)
			apiStampsSID := apiStamps.Group("/:stampID")
			{
				apiStampsSID.GET("", NotImplemented)
				apiStampsSID.DELETE("", NotImplemented)
			}
		}
		apiStampPalettes := api.Group("/stamp-palettes")
		{
			apiStampPalettes.GET("", NotImplemented)
			apiStampPalettes.POST("", NotImplemented)
			apiStampPalettesPID := apiStampPalettes.Group("/:paletteID")
			{
				apiStampPalettesPID.GET("", NotImplemented)
				apiStampPalettesPID.PATCH("", NotImplemented)
				apiStampPalettesPID.DELETE("", NotImplemented)
				apiStampPalettesPID.PUT("/stamps", NotImplemented)
			}
		}
		apiWebhooks := api.Group("/webhooks")
		{
			apiWebhooks.GET("", NotImplemented)
			apiWebhooks.POST("", NotImplemented)
			apiWebhooksWID := apiWebhooks.Group("/:webhookID")
			{
				apiWebhooksWID.GET("", NotImplemented)
				apiWebhooksWID.PATCH("", NotImplemented)
				apiWebhooksWID.DELETE("", NotImplemented)
				apiWebhooksWID.GET("/icon", NotImplemented)
				apiWebhooksWID.PUT("/icon", NotImplemented)
				apiWebhooksWID.GET("/messages", NotImplemented)
			}
		}
		apiGroups := api.Group("/groups")
		{
			apiGroups.GET("", NotImplemented)
			apiGroups.POST("", NotImplemented)
			apiGroupsGID := apiGroups.Group("/:groupID")
			{
				apiGroupsGID.GET("", NotImplemented)
				apiGroupsGID.PATCH("", NotImplemented)
				apiGroupsGID.DELETE("", NotImplemented)
				apiGroupsGIDMembers := apiGroupsGID.Group("/members")
				{
					apiGroupsGIDMembers.GET("", NotImplemented)
					apiGroupsGIDMembers.POST("", NotImplemented)
					apiGroupsGIDMembers.PUT("", NotImplemented)
					apiGroupsGIDMembersUID := apiGroupsGIDMembers.Group("/:userID")
					{
						apiGroupsGIDMembersUID.PATCH("", NotImplemented)
						apiGroupsGIDMembersUID.DELETE("", NotImplemented)
					}
				}
			}
		}
		apiActivity := api.Group("/activity")
		{
			apiActivity.GET("/timelines", NotImplemented)
			apiActivity.GET("/onlines", h.GetOnlineUsers)
		}
		apiClients := api.Group("/clients")
		{
			apiClients.GET("", NotImplemented)
			apiClients.POST("", NotImplemented)
			apiClientsCID := apiClients.Group("/:clientID")
			{
				apiClientsCID.GET("", NotImplemented)
				apiClientsCID.PATCH("", NotImplemented)
				apiClientsCID.DELETE("", NotImplemented)
			}
		}
		apiBots := api.Group("/bots")
		{
			apiBots.GET("", NotImplemented)
			apiBots.POST("", NotImplemented)
			apiBotsBID := apiBots.Group("/:botID")
			{
				apiBotsBID.GET("", NotImplemented)
				apiBotsBID.PATCH("", NotImplemented)
				apiBotsBID.DELETE("", NotImplemented)
				apiBotsBID.GET("/icon", NotImplemented)
				apiBotsBID.PUT("/icon", NotImplemented)
				apiBotsBID.GET("/logs", NotImplemented)
				apiBotsBIDActions := apiBotsBID.Group("/actions")
				{
					apiBotsBIDActions.POST("/activate", NotImplemented)
					apiBotsBIDActions.POST("/inactivate", NotImplemented)
					apiBotsBIDActions.POST("/reissue", NotImplemented)
					apiBotsBIDActions.POST("/join", NotImplemented)
					apiBotsBIDActions.POST("/leave", NotImplemented)
				}
			}
		}
		apiWebRTC := api.Group("/webrtc")
		{
			apiWebRTC.GET("/state", NotImplemented)
			apiWebRTC.POST("/authenticate", NotImplemented)
		}
		apiClipFolders := api.Group("/clip-folders")
		{
			apiClipFolders.GET("", NotImplemented)
			apiClipFolders.POST("", NotImplemented)
			apiClipFoldersFID := apiClipFolders.Group("/:folderID")
			{
				apiClipFoldersFID.GET("", NotImplemented)
				apiClipFoldersFID.PATCH("", NotImplemented)
				apiClipFoldersFID.DELETE("", NotImplemented)
				apiClipFoldersFIDMessages := apiClipFoldersFID.Group("/messages")
				{
					apiClipFoldersFIDMessages.GET("", NotImplemented)
					apiClipFoldersFIDMessages.POST("", NotImplemented)
					apiClipFoldersFIDMessages.DELETE("/:messageID", NotImplemented)
				}
			}
		}
		api.GET("/ws", echo.WrapHandler(h.WS), requires(permission.ConnectNotificationStream))
	}

	apiNoAuth := e.Group("/v3")
	{
		apiNoAuth.GET("/version", h.GetVersion)
		apiNoAuth.POST("/login", NotImplemented)
		apiNoAuth.POST("/logout", NotImplemented)
		apiNoAuth.POST("/webhooks/:webhookID", NotImplemented)
		apiNoAuthPublic := apiNoAuth.Group("/public")
		{
			apiNoAuthPublic.GET("/icon/:username", NotImplemented)
		}
	}
}
