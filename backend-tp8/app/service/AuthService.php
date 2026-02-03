<?php

declare(strict_types=1);

namespace app\service;

use app\model\AdminUser;
use app\model\PointsAccount;
use app\model\PointsRule;
use app\model\User;
use app\model\UserProfile;
use app\model\PointsLog;
use think\facade\Db;

class AuthService
{
    public static function adminLogin(string $username, string $password): array
    {
        $admin = AdminUser::where('username', $username)->find();
        if (!$admin || !password_verify($password, $admin->password_hash)) {
            throw new \RuntimeException('用户名或密码错误');
        }
        if ((int) $admin->status !== 1) {
            throw new \RuntimeException('账号已禁用');
        }

        $token = JwtService::encode([
            'role' => 'admin',
            'admin_id' => $admin->id,
            'username' => $admin->username,
        ]);

        return [
            'token' => $token,
            'admin' => [
                'id' => $admin->id,
                'username' => $admin->username,
            ],
        ];
    }

    public static function userRegister(string $username, string $password): array
    {
        return Db::transaction(function () use ($username, $password) {
            $exists = User::where('username', $username)->find();
            if ($exists) {
                throw new \RuntimeException('用户名已存在');
            }

            $hash = password_hash($password, PASSWORD_BCRYPT);
            $user = User::create([
                'username' => $username,
                'password_hash' => $hash,
                'status' => 1,
            ]);

            UserProfile::create([
                'user_id' => $user->id,
            ]);

            $rule = PointsRule::order('id', 'desc')->find();
            $bonus = $rule ? (int) $rule->register_bonus_points : 0;

            PointsAccount::create([
                'user_id' => $user->id,
                'available_points' => $bonus,
                'frozen_points' => 0,
                'withdrawn_points' => 0,
            ]);

            if ($bonus > 0) {
                PointsLog::create([
                    'user_id' => $user->id,
                    'change_points' => $bonus,
                    'type' => 'register_bonus',
                    'ref_id' => null,
                    'remark' => '注册赠送积分',
                ]);
            }

            $token = JwtService::encode([
                'role' => 'user',
                'user_id' => $user->id,
                'username' => $user->username,
            ]);

            return [
                'token' => $token,
                'user' => [
                    'id' => $user->id,
                    'username' => $user->username,
                ],
            ];
        });
    }

    public static function userLogin(string $username, string $password): array
    {
        $user = User::where('username', $username)->find();
        if (!$user || !password_verify($password, $user->password_hash)) {
            throw new \RuntimeException('用户名或密码错误');
        }
        if ((int) $user->status !== 1) {
            throw new \RuntimeException('账号已禁用');
        }

        $token = JwtService::encode([
            'role' => 'user',
            'user_id' => $user->id,
            'username' => $user->username,
        ]);

        return [
            'token' => $token,
            'user' => [
                'id' => $user->id,
                'username' => $user->username,
            ],
        ];
    }
}
