
CREATE TABLE public.todo (
  id serial NOT NULL,
  title varchar(255) NOT NULL,
  detail text NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz NULL,
  CONSTRAINT "todo_id" PRIMARY KEY ("id")
)
WITH (
	OIDS=FALSE
);