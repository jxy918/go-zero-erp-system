import re
with open('d:/work/go/src/go-zero-erp/admin/internal/types/types.go', 'r', encoding='utf-8') as f:
    content = f.read()

old = '''type ListProductCategoryRequest struct {
	Page     int `form:"page,optional"`
	PageSize int `form:"page_size,optional"`
}'''

new = '''type ListProductCategoryRequest struct {
	Page     int    `form:"page,optional"`
	PageSize int    `form:"page_size,optional"`
	Name     string `form:"name,optional"`
}'''

content = content.replace(old, new)

with open('d:/work/go/src/go-zero-erp/admin/internal/types/types.go', 'w', encoding='utf-8') as f:
    f.write(content)

print('Done')
