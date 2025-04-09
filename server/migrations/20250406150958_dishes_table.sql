-- +goose Up
-- +goose StatementBegin
create table if not exist dishes(
    id uuid primary key,
    name text not null,
    recipe jsonb not null,
    ingredients jsonb not null,
    nutrition jsonb not null,
    meal_type text not null,
    cooking_time int32 not null,
    servings int32 not null
)
CREATE INDEX idx_dishes_meal_type ON dishes(meal_type);
CREATE INDEX idx_dishes_nutrition ON dishes USING GIN(nutrition);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index idx_dishes_meal_type
drop index idx_dishes_nutrition
drop table dishes
-- +goose StatementEnd
