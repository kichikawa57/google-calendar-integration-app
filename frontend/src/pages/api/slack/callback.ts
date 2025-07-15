// pages/api/slack/callback.ts

import type { NextApiRequest, NextApiResponse } from "next";

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  const code = req.query.code as string;

  if (!code) {
    return res.status(400).send("Code not found");
  }

  const client_id = process.env.SLACK_CLIENT_ID;
  const client_secret = process.env.SLACK_CLIENT_SECRET;

  // https://api.slack.com/methods/oauth.v2.access を参考にしてください。
  const tokenResponse = await fetch("https://slack.com/api/oauth.v2.access", {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body: new URLSearchParams({
      code,
      client_id: client_id || "",
      client_secret: client_secret || "",
    }),
  });

  const data = await tokenResponse.json();

  if (!data.ok) {
    return res.status(500).json({ error: data.error });
  }

  // 🔽 Webhook URL がここに含まれる
  const webhookUrl = data.incoming_webhook?.url;

  // 必要に応じてDB保存など
  console.log("Webhook URL:", webhookUrl);

  // 任意のページにリダイレクト
  return res.redirect("/");
}
