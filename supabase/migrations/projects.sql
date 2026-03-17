-- 1. Drop the old table if it exists to start clean
DROP TABLE IF EXISTS public.projects;

-- 2. Create the simplified projects table
CREATE TABLE public.projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    package_name TEXT NOT NULL,
    platform TEXT DEFAULT 'android',
    color TEXT DEFAULT '#2e60ff',
    icon_url TEXT,
    version_code INTEGER DEFAULT 1,
    version_name TEXT DEFAULT '1.0.0',
    
    -- Workspace and code storage
    workspace_xml TEXT,  -- Blockly workspace XML
    generated_code TEXT,  -- Generated React Native code
    
    -- Status and build tracking
    status TEXT DEFAULT 'draft',
    last_build_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 3. Enable RLS
ALTER TABLE public.projects ENABLE ROW LEVEL SECURITY;

-- 4. Policy: Users can only manage their own projects
CREATE POLICY "Users can manage their own projects" 
ON public.projects FOR ALL 
USING (auth.uid() = user_id)
WITH CHECK (auth.uid() = user_id);

-- 5. Enable Realtime
ALTER PUBLICATION supabase_realtime ADD TABLE public.projects;
