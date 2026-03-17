-- 1. Create the table for global user keystores
CREATE TABLE IF NOT EXISTS public.user_keystores (
    user_id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
    keystore_name TEXT NOT NULL, -- Will be user_id.keystore
    alias TEXT NOT NULL DEFAULT 'rena_key',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. Enable RLS
ALTER TABLE public.user_keystores ENABLE ROW LEVEL SECURITY;

-- 3. Policy: Users can only see/manage their own keystore metadata
CREATE POLICY "Users can manage their own keystore" 
ON public.user_keystores FOR ALL 
USING (auth.uid() = user_id)
WITH CHECK (auth.uid() = user_id);
