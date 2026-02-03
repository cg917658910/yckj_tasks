<?php

declare(strict_types=1);

namespace app\service;

use Firebase\JWT\JWT;
use Firebase\JWT\Key;
use think\facade\Config;

class JwtService
{
    public static function encode(array $payload): string
    {
        $config = Config::get('jwt');
        $time = time();

        $claims = array_merge($payload, [
            'iss' => $config['issuer'],
            'iat' => $time,
            'exp' => $time + (int) $config['ttl'],
        ]);

        return JWT::encode($claims, $config['secret'], 'HS256');
    }

    public static function decode(string $token): array
    {
        $config = Config::get('jwt');
        $decoded = JWT::decode($token, new Key($config['secret'], 'HS256'));
        return (array) $decoded;
    }
}
