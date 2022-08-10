-- +goose Up
-- +goose StatementBegin
INSERT INTO public.users (name, password, email, full_name, created_at)
VALUES
    ('Piter','123','piter@email.com','Piter Parker',1659447420),
    ('Sara','321','sara@email.com','Sara Conor',1659447430),
    ('Tony','456','tony@email.com','Tony Stark',1659447440),
    ('Berta','654','berta@email.com','Big Berta',1659447450);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM public.users
WHERE name IN ('Piter','Sara','Tony','Berta');
-- +goose StatementEnd
