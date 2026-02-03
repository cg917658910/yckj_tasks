<?php

declare(strict_types=1);

namespace app\service;

use app\model\PointsAccount;
use app\model\PointsLog;
use app\model\PointsRule;
use think\facade\Db;

class PointsService
{
    public static function getRules(): PointsRule
    {
        $rule = PointsRule::order('id', 'desc')->find();
        if (!$rule) {
            $rule = PointsRule::create([
                'exchange_rate' => 10,
                'min_withdraw_amount' => 10.00,
                'register_bonus_points' => 10,
            ]);
        }
        return $rule;
    }

    public static function updateRules(int $adminId, array $data): PointsRule
    {
        return Db::transaction(function () use ($adminId, $data) {
            $rule = self::getRules();
            $old = $rule->toArray();
            $rule->save($data);

            Db::name('points_rule_logs')->insert([
                'rule_id' => $rule->id,
                'old_value' => json_encode($old, JSON_UNESCAPED_UNICODE),
                'new_value' => json_encode($rule->toArray(), JSON_UNESCAPED_UNICODE),
                'admin_id' => $adminId,
                'changed_at' => date('Y-m-d H:i:s'),
            ]);

            return $rule;
        });
    }

    public static function addAvailable(int $userId, int $points, string $type, ?int $refId, string $remark): void
    {
        Db::transaction(function () use ($userId, $points, $type, $refId, $remark) {
            $account = PointsAccount::where('user_id', $userId)->lock(true)->find();
            if (!$account) {
                $account = PointsAccount::create([
                    'user_id' => $userId,
                    'available_points' => 0,
                    'frozen_points' => 0,
                    'withdrawn_points' => 0,
                ]);
            }
            $account->available_points += $points;
            $account->save();

            PointsLog::create([
                'user_id' => $userId,
                'change_points' => $points,
                'type' => $type,
                'ref_id' => $refId,
                'remark' => $remark,
            ]);
        });
    }

    public static function freezeForWithdraw(int $userId, int $points, string $remark): void
    {
        Db::transaction(function () use ($userId, $points, $remark) {
            $account = PointsAccount::where('user_id', $userId)->lock(true)->findOrFail();
            if ($account->available_points < $points) {
                throw new \RuntimeException('积分不足');
            }
            $account->available_points -= $points;
            $account->frozen_points += $points;
            $account->save();

            PointsLog::create([
                'user_id' => $userId,
                'change_points' => -$points,
                'type' => 'withdraw_freeze',
                'ref_id' => null,
                'remark' => $remark,
            ]);
        });
    }

    public static function unfreezeToAvailable(int $userId, int $points, string $remark): void
    {
        Db::transaction(function () use ($userId, $points, $remark) {
            $account = PointsAccount::where('user_id', $userId)->lock(true)->findOrFail();
            if ($account->frozen_points < $points) {
                throw new \RuntimeException('冻结积分不足');
            }
            $account->frozen_points -= $points;
            $account->available_points += $points;
            $account->save();

            PointsLog::create([
                'user_id' => $userId,
                'change_points' => $points,
                'type' => 'withdraw_unfreeze',
                'ref_id' => null,
                'remark' => $remark,
            ]);
        });
    }

    public static function deductFrozenToWithdrawn(int $userId, int $points, string $remark): void
    {
        Db::transaction(function () use ($userId, $points, $remark) {
            $account = PointsAccount::where('user_id', $userId)->lock(true)->findOrFail();
            if ($account->frozen_points < $points) {
                throw new \RuntimeException('冻结积分不足');
            }
            $account->frozen_points -= $points;
            $account->withdrawn_points += $points;
            $account->save();

            PointsLog::create([
                'user_id' => $userId,
                'change_points' => -$points,
                'type' => 'withdraw_complete',
                'ref_id' => null,
                'remark' => $remark,
            ]);
        });
    }
}
