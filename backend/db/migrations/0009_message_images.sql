-- +goose Up
-- A message may now carry a private image (proof of registration / handover)
-- instead of, or alongside, text. Body becomes optional; a message must still
-- carry text or an image.
ALTER TABLE messages ADD COLUMN image_key text;
ALTER TABLE messages ALTER COLUMN body DROP NOT NULL;
ALTER TABLE messages DROP CONSTRAINT messages_body_length;
ALTER TABLE messages ADD CONSTRAINT messages_content_check CHECK (
    (body IS NULL OR char_length(body) BETWEEN 1 AND 4000)
    AND (body IS NOT NULL OR image_key IS NOT NULL)
);

-- +goose Down
ALTER TABLE messages DROP CONSTRAINT messages_content_check;
ALTER TABLE messages ALTER COLUMN body SET NOT NULL;
ALTER TABLE messages
    ADD CONSTRAINT messages_body_length CHECK (char_length(body) BETWEEN 1 AND 4000);
ALTER TABLE messages DROP COLUMN image_key;
