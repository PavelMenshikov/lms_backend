-- +goose Up
-- +goose StatementBegin
ALTER TABLE lessons ADD COLUMN IF NOT EXISTS content JSONB DEFAULT '[]'::jsonb;

UPDATE lessons 
SET content = jsonb_build_array(
    jsonb_build_object('type', 'text', 'value', content_text)
)
WHERE content_text IS NOT NULL AND content = '[]'::jsonb;
-- +goose StatementEnd

-- +goose Down
ALTER TABLE lessons DROP COLUMN IF EXISTS content;