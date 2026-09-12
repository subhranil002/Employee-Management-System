import react from "@vitejs/plugin-react";
import { defineConfig, loadEnv } from "vite";
import tailwindcss from "@tailwindcss/vite";
import { SecretsManagerClient, GetSecretValueCommand } from "@aws-sdk/client-secrets-manager";

// https://vite.dev/config/
export default defineConfig(async ({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');
  const secretArn = env.FRONTEND_SECRET_ARN;
  let apiUrl = env.VITE_API_BASE_URL;

  // AWS: load configuration from Secrets Manager
  if (secretArn) {
    try {
      const client = new SecretsManagerClient({ region: env.AWS_REGION || "ap-south-1" });
      const response = await client.send(new GetSecretValueCommand({ SecretId: secretArn }));
      
      if (response.SecretString) {
        const secret = JSON.parse(response.SecretString);
        
        if (secret.VITE_API_BASE_URL) {
          apiUrl = secret.VITE_API_BASE_URL;
        }

        // Map other VITE_ prefixed secrets to environment variables
        for (const [key, value] of Object.entries(secret)) {
          if (key.startsWith("VITE_")) {
            process.env[key] = value as string;
          }
        }
      }
    } catch (err) {
      console.error("Failed to load secrets from AWS Secrets Manager:", err);
    }
  }

  return {
    plugins: [tailwindcss(), react()],
    define: apiUrl ? {
      'import.meta.env.VITE_API_BASE_URL': JSON.stringify(apiUrl),
    } : undefined,
  };
});
