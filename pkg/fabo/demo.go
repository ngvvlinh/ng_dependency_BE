package fabo

import (
	"context"

	"o.o/api/top/int/fabo"
	"o.o/backend/com/fabo/pkg/fbclient"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/authorize/session"
)

type DemoService struct {
	session.Session

	FBClient *fbclient.FbClient
}

func (s *DemoService) Clone() fabo.DemoService {
	res := *s
	return &res
}

func (s *DemoService) ListLiveVideos(
	ctx context.Context, req *fabo.ListLiveVideosRequest,
) (*fabo.ListLiveVideosResponse, error) {
	listLiveVideosReq := &fbclient.ListLiveVideosRequest{
		AccessToken: req.Token,
	}
	listLiveVideosResp, err := s.FBClient.CallAPIListLiveVideos(listLiveVideosReq)
	if err != nil {
		return nil, cm.Errorf(cm.ExternalServiceError, err, "Error from Facebook: %s", err.Error())
	}

	var videos []*fabo.LiveVideoUser

	if listLiveVideosResp.LiveVideos != nil && len(listLiveVideosResp.LiveVideos.Data) > 0 {
		videosData := listLiveVideosResp.LiveVideos.Data

		for _, videoData := range videosData {
			liveVideoUser := &fabo.LiveVideoUser{
				ID:           videoData.ID,
				Title:        videoData.Title,
				Description:  videoData.Description,
				PermalinkURL: videoData.PermalinkURL,
				EmbedHTML:    videoData.EmbedHTML,
				CreatedTime:  videoData.CreationTime.ToTime(),
			}
			if videoData.From != nil {
				liveVideoUser.From = &fabo.FbObjectFrom{
					ID:    videoData.From.ID,
					Name:  videoData.From.Name,
					Email: videoData.From.Email,
				}
			}

			if videoData.Comments != nil {
				var comments []*fabo.LiveVideoComment
				for _, commentData := range videoData.Comments.Data {
					comments = append(comments, &fabo.LiveVideoComment{
						CreatedTime: commentData.CreatedTime.ToTime(),
						ID:          commentData.ID,
						Message:     commentData.Message,
					})
				}

				liveVideoUser.Comments = comments
			}

			if videoData.Video != nil {
				liveVideoUser.Video = &fabo.LiveVideoVideo{
					ID:      videoData.Video.ID,
					Picture: videoData.Video.Picture,
					Source:  videoData.Video.Source,
				}
			}

			videos = append(videos, liveVideoUser)
		}

	}

	return &fabo.ListLiveVideosResponse{
		Videos: videos,
	}, nil
}

func (s *DemoService) ListFeeds(
	ctx context.Context, req *fabo.ListFeedsRequest,
) (*fabo.ListFeedsResponse, error) {
	listFeedsWithCommentsReq := &fbclient.ListFeedsWithCommentsRequest{
		AccessToken: req.Token,
	}
	listFeedsWithCommentsResp, err := s.FBClient.CallAPIListFeedsWithComments(listFeedsWithCommentsReq)
	if err != nil {
		return nil, cm.Errorf(cm.ExternalServiceError, err, "Error from Facebook: %s", err.Error())
	}

	var feeds []*fabo.PostWithComments
	if listFeedsWithCommentsResp.Feeds != nil && len(listFeedsWithCommentsResp.Feeds.Data) > 0 {
		feedsData := listFeedsWithCommentsResp.Feeds.Data

		for _, feedData := range feedsData {
			feed := &fabo.PostWithComments{
				Post: fabo.Post{
					ID: feedData.ID,

					FullPicture:  feedData.FullPicture,
					Icon:         feedData.Icon,
					IsExpired:    feedData.IsExpired,
					IsHidden:     feedData.IsHidden,
					IsPopular:    feedData.IsPopular,
					IsPublished:  feedData.IsPublished,
					Message:      feedData.Message,
					Story:        feedData.Story,
					PermalinkURL: feedData.PermalinkURL,
					StatusType:   feedData.StatusType,
					Picture:      feedData.Picture,
					CreatedTime:  feedData.CreatedTime.ToTime(),
					UpdatedTime:  feedData.UpdatedTime.ToTime(),
				},
			}

			if feedData.From != nil {
				feed.Post.From = &fabo.FbObjectFrom{
					ID:    feedData.From.ID,
					Name:  feedData.From.Name,
					Email: feedData.From.Email,
				}
			}

			if feedData.Attachments != nil {
				var attachments []*fabo.PostAttachment
				for _, feedAttachment := range feedData.Attachments.Data {
					attachment := &fabo.PostAttachment{
						MediaType:      feedAttachment.MediaType,
						Type:           feedAttachment.Type,
						SubAttachments: nil,
					}

					if feedAttachment.Media != nil && feedAttachment.Media.Image != nil {
						attachment.Media = &fabo.MediaPostAttachment{
							Image: &fabo.ImageMediaPostAttachment{
								Height: feedAttachment.Media.Image.Height,
								Width:  feedAttachment.Media.Image.Width,
								Src:    feedAttachment.Media.Image.Src,
							},
						}
					}

					if feedAttachment.SubAttachments != nil {
						var subAttachments []*fabo.SubAttachment
						for _, feedSubAttachment := range feedAttachment.SubAttachments.Data {
							subAttachment := &fabo.SubAttachment{
								Type: feedSubAttachment.Type,
								URL:  feedSubAttachment.URL,
							}

							if feedSubAttachment.Media != nil && feedSubAttachment.Media.Image != nil {
								subAttachment.Media = &fabo.MediaDataSubAttachment{
									Height: feedSubAttachment.Media.Image.Height,
									Width:  feedSubAttachment.Media.Image.Width,
									Src:    feedSubAttachment.Media.Image.Src,
								}
							}

							if feedSubAttachment.Target != nil {
								subAttachment.Target = &fabo.TargetDataSubAttachment{
									ID:  feedSubAttachment.Target.ID,
									URL: feedSubAttachment.Target.URL,
								}
							}

							subAttachments = append(subAttachments, subAttachment)
						}

						attachment.SubAttachments = subAttachments
					}

					attachments = append(attachments, attachment)
				}

				feed.Attachments = attachments
			}

			if feedData.Comments != nil {
				var comments []*fabo.PostComment
				for _, feedComment := range feedData.Comments.CommentData {
					comments = append(comments, &fabo.PostComment{
						CreatedTime: feedComment.CreatedTime.ToTime(),
						ID:          feedComment.ID,
						Message:     feedComment.Message,
					})
				}
				feed.Comments = comments
			}

			feeds = append(feeds, feed)
		}
	}

	return &fabo.ListFeedsResponse{
		Feeds: feeds,
	}, nil
}
