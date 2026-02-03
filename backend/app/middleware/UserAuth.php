<?php

declare(strict_types=1);

namespace app\middleware;

use app\service\JwtService;
use think\Request;
use think\Response;

class UserAuth
{
    public function handle(Request $request, \Closure $next): Response
    {
        $token = $this->getBearerToken($request);
        if ($token === null) {
            return json_error('未登录', 401, [], 401);
        }

        try {
            $payload = JwtService::decode($token);
        } catch (\Throwable $e) {
            return json_error('登录已失效', 401, [], 401);
        }

        if (($payload['role'] ?? '') !== 'user') {
            return json_error('无权限', 403, [], 403);
        }

        $request->user = $payload;

        return $next($request);
    }

    private function getBearerToken(Request $request): ?string
    {
        $header = $request->header('Authorization');
        if (!$header) {
            return null;
        }
        if (stripos($header, 'Bearer ') !== 0) {
            return null;
        }
        return trim(substr($header, 7));
    }
}
