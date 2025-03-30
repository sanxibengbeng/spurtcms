-- Check if the table exists
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'tbl_storage_types') THEN
        -- Create the table if it doesn't exist
        CREATE TABLE tbl_storage_types (
            id SERIAL PRIMARY KEY,
            local VARCHAR(255),
            aws JSONB,
            azure JSONB,
            drive JSONB,
            selected_type VARCHAR(50)
        );
        
        -- Insert default record with local storage configuration
        INSERT INTO tbl_storage_types (id, local, selected_type)
        VALUES (1, 'storage', 'local');
    ELSE
        -- Update the existing record to use local storage
        UPDATE tbl_storage_types
        SET local = 'storage', selected_type = 'local'
        WHERE id = 1;
        
        -- If no record exists, insert one
        INSERT INTO tbl_storage_types (id, local, selected_type)
        SELECT 1, 'storage', 'local'
        WHERE NOT EXISTS (SELECT 1 FROM tbl_storage_types WHERE id = 1);
    END IF;
END
$$;
