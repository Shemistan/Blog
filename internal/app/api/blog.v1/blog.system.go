package blog_v1

import (
	blog_system "github.com/Shemistan/Blog/internal/app/service/blog.system"
	pb "github.com/Shemistan/Blog/pkg/blog.v1"
)

type Blog struct {
	// Нужно, что бы приложение не падало в панике, если какой-то АПИ еще не реализован.
	pb.UnimplementedBlogV1Server

	BlogService blog_system.IBlogSystemService
}
