ALTER TABLE auth_providers
ADD CONSTRAINT auth_providers_pkey PRIMARY KEY (provider, provider_user_id);