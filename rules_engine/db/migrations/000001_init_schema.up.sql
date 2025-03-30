CREATE TABLE resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    http_method TEXT NOT NULL CHECK (http_method IN ('GET', 'POST', 'PUT', 'DELETE', 'PATCH')),
    url TEXT NOT NULL,
    creator_id UUID NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    attack_type TEXT NOT NULL CHECK (attack_type IN ('xss', 'csrf', 'sqli')),
    action_type TEXT NOT NULL CHECK (action_type IN ('block', 'sanitize', 'escape')),
    is_active BOOLEAN DEFAULT TRUE,
    creator_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE resource_rule (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_id UUID NOT NULL REFERENCES rules(id) ON DELETE CASCADE,
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (rule_id, resource_id)
);

CREATE TABLE ip_lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ip CIDR NOT NULL,
    list_type TEXT NOT NULL CHECK (list_type IN ('whitelist', 'blacklist')),
    creator_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (ip, list_type)
);

CREATE TABLE resource_ip_list (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    ip_list_id UUID NOT NULL REFERENCES ip_lists(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);
