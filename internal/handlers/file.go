package handlers

import (
	"bytes"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	fs "mmplat/internal/filesystem"
	"mmplat/internal/templates"
	"mmplat/internal/util"
	"os"
	"strconv"
	"strings"
)

// Index basically, handler for inode/directory
// URL/{ID}/{ID}
// TODO implement that thing above + plus breadcrumbs(or .. in workdir tree)

const assetDir = "../../assets"

func (h *Handler) Index(ctx *fasthttp.RequestCtx) {
	h.index(nil, ctx)
}

// Item file handler
func (h *Handler) Item(ctx *fasthttp.RequestCtx) {
	tree := h.fs.Tree()
	id, _ := strconv.Atoi(ctx.UserValue("item").(string))
	if node := tree.FindId(fs.Id(id)); node != nil {
		switch node.Type() {
		case fs.Dir:
			h.index(node, ctx)
		case fs.File:
			h.render(node, ctx)
		default:
		}
	} else {
		util.SetTitle(ctx, "File not found.")
		h.NotFound(ctx)
	}
}

func (h *Handler) index(node *fs.Node, ctx *fasthttp.RequestCtx) {
	var files []map[string]string
	var searchNode *fs.Node
	searchNode = node
	if node == nil {
		searchNode = h.fs.Tree().Root()
	}
	for _, n := range searchNode.Children() {
		files = append(files, util.PrepareTemplateItem(n.Id(), n.Item()))
	}

	p := &templates.FileListPage{Files: files}
	templates.WritePageTemplate(ctx, p)
}

func (h *Handler) render(node *fs.Node, ctx *fasthttp.RequestCtx) {
	reader, err := node.Item().Reader()
	// TODO non beautiful switch for common types
	if strings.Contains(node.Item().Metadata(), "video") {
		if err != nil {
			fmt.Fprintf(ctx, "Error reading file")
		}
		file := util.PrepareTemplateItem(node.Id(), node.Item())
		var p templates.Page
		// redir no js
		if ctx.QueryArgs().Has("nojs") {
			p = &templates.RenderNoJs{File: file}
		} else {
			p = &templates.Render{File: file}
		}
		templates.WritePageTemplate(ctx, p)
	} else {
		ctx.SetBodyStream(reader, -1)
		metadata := util.ExtToMetadata(node.Item())
		ctx.Response.Header.Set("Content-Disposition",
			fmt.Sprintf(`inline; filename="%s`, node.Item().Name()),
		)
		if metadata != "" {
			ctx.SetContentType(metadata)
		} else {
			ctx.SetContentType("application/octet-stream")
		}
	}
}

func (h *Handler) FaviconPieceOfShit(ctx *fasthttp.RequestCtx) {
	asset := assetDir + "/" + "img/favicon.ico"
	file, err := os.Open(asset)
	if err != nil {
		h.NotFound(ctx)
	}
	ctx.SetBodyStream(file, -1)
}

func (h *Handler) Asset(ctx *fasthttp.RequestCtx) {
	asset := assetDir + "/" + util.ValidateUserInput(ctx.UserValue("asset").(string))
	var file *os.File
	var err error
	file, err = os.Open(asset)
	if strings.HasSuffix(asset, "js") {
		ctx.SetContentType("text/javascript")
	} else if strings.HasSuffix(asset, "js") {
		ctx.SetContentType("text/css")
	}
	if err != nil {
		h.NotFound(ctx)
		return
	}
	ctx.SetBodyStream(file, -1)
}

// Stream separte function to stream using filename as asource
func (h *Handler) Stream(ctx *fasthttp.RequestCtx) {
	tree := h.fs.Tree()
	name := util.ValidateUserInput(ctx.UserValue("name").(string))
	if node := tree.FindBy(func(needle *fs.Node) *fs.Node {
		if needle.Item().Name() == name {
			return needle
		}
		return nil
	}); node != nil {
		// 1) why 1st request fails miserably
		// 2nd, check the bytes part and offset f-on response to browser
		if rng := util.CheckRange(ctx); rng > -1 {
			data, _ := ContentChunk(node, rng)
			ctx.SetContentType(node.Item().Metadata())
			ctx.Response.Header.Add("Content-Length", strconv.Itoa(int(node.Item().Size())))
			ctx.Response.Header.Add("Content-Range", fmt.Sprintf("bytes %s-%s/%s", strconv.Itoa(rng),
				strconv.Itoa(int(node.Item().Size())-1),
				strconv.Itoa(int(node.Item().Size()))),
			)
			ctx.SetStatusCode(fasthttp.StatusPartialContent)
			ctx.SetBody(data)
		}
	}
}

// ContentChunk lets guess content chunk, by 5-10 sec cache universaly
// (or just return 20mb chunk)
func ContentChunk(node *fs.Node, offset int) ([]byte, int) {
	file, err := node.Item().Reader()
	buf := new(bytes.Buffer)
	temp := make([]byte, 2*12)
	if err == nil {
		defer file.Close()
		_, err := file.Seek(int64(offset), 0)
		if err == nil {
			var l int
			for {
				rd, err := file.Read(temp)
				l += rd
				buf.Write(temp)
				if l >= 20*1024*1024 || err == io.EOF {
					break
				}
			}
			return buf.Bytes(), l
		}
	}
	return nil, 0
}
