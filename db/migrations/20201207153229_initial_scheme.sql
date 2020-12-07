
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

--
-- Name: users; Type: TABLE; Schema: public; Owner: master
--

CREATE TABLE IF NOT EXISTS public.users (
    id SERIAL PRIMARY KEY,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    provider character varying(255) DEFAULT NULL::character varying,
    oauth_token character varying(255) DEFAULT NULL::character varying,
    uuid bigint,
    user_name character varying(255) DEFAULT NULL::character varying,
    avatar_url character varying(255) DEFAULT NULL::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

--
-- Name: email_on_users; Type: INDEX; Schema: public; Owner: master
--

CREATE UNIQUE INDEX email_on_users ON public.users USING btree (email);


--
-- Name: inquiries; Type: TABLE; Schema: public; Owner: master
--

CREATE TABLE IF NOT EXISTS public.inquiries (
    id SERIAL PRIMARY KEY,
    email character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    message text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

--
-- Name: repositories; Type: TABLE; Schema: public; Owner: master
--

CREATE TABLE IF NOT EXISTS public.repositories (
    id SERIAL PRIMARY KEY,
    repository_id bigint NOT NULL,
    owner character varying(255) DEFAULT NULL::character varying,
    name character varying(255) DEFAULT NULL::character varying,
    webhook_key character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

--
-- Name: repository_id_on_repositories; Type: INDEX; Schema: public; Owner: master
--

CREATE UNIQUE INDEX repository_id_on_repositories ON public.repositories USING btree (repository_id);


--
-- Name: projects; Type: TABLE; Schema: public; Owner: master
--

CREATE TABLE IF NOT EXISTS public.projects (
    id SERIAL PRIMARY KEY,
    user_id integer REFERENCES users(id) NOT NULL,
    repository_id integer REFERENCES repositories(id),
    title character varying(255) NOT NULL,
    description character varying(255) NOT NULL,
    show_issues boolean DEFAULT true NOT NULL,
    show_pull_requests boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: title_and_user_id_on_projects; Type: INDEX; Schema: public; Owner: master
--

CREATE UNIQUE INDEX title_and_user_id_on_projects ON public.projects USING btree (title, user_id);


--
-- Name: list_options; Type: TABLE; Schema: public; Owner: master
--

CREATE TABLE public.list_options (
    id SERIAL PRIMARY KEY,
    action character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

--
-- Name: action_on_list_options; Type: INDEX; Schema: public; Owner: master
--

CREATE UNIQUE INDEX action_on_list_options ON public.list_options USING btree (action);

--
-- Name: lists; Type: TABLE; Schema: public; Owner: master
--

CREATE TABLE IF NOT EXISTS public.lists (
    id SERIAL PRIMARY KEY,
    project_id integer REFERENCES projects(id) NOT NULL,
    user_id integer REFERENCES users(id) NOT NULL,
    list_option_id integer REFERENCES list_options(id),
    title character varying(255) DEFAULT NULL::character varying,
    color character varying(255) DEFAULT NULL::character varying,
    is_hidden boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: tasks; Type: TABLE; Schema: public; Owner: master
--

CREATE TABLE IF NOT EXISTS public.tasks (
    id SERIAL PRIMARY KEY,
    user_id integer REFERENCES users(id) NOT NULL,
    project_id integer REFERENCES projects(id) NOT NULL,
    list_id integer REFERENCES lists(id) NOT NULL,
    title character varying(255) NOT NULL,
    display_index integer NOT NULL,
    issue_number integer,
    description text NOT NULL,
    pull_request boolean DEFAULT false NOT NULL,
    html_url character varying(255) DEFAULT NULL::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

--
-- Name: project_id_and_issue_number_on_tasks; Type: INDEX; Schema: public; Owner: master
--

CREATE UNIQUE INDEX project_id_and_issue_number_on_tasks ON public.tasks USING btree (project_id, issue_number);




-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP INDEX project_id_and_issue_number_on_tasks;
DROP TABLE public.tasks;
DROP TABLE public.lists;
DROP INDEX action_on_list_options;
DROP TABLE public.list_options;
DROP INDEX title_and_user_id_on_projects;
DROP TABLE public.projects;
DROP INDEX repository_id_on_repositories;
DROP TABLE public.repositories;
DROP TABLE public.inquiries;
DROP INDEX email_on_users;
DROP TABLE public.users;
