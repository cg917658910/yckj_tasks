<?php

declare(strict_types=1);

namespace app\service;

use app\model\PointsRule;
use app\model\UserProfile;
use app\model\Withdrawal;
use think\facade\Db;

class WithdrawalService
{
    public const STATUS_PENDING = 1;
    public const STATUS_PAID = 2;
    public const STATUS_REJECTED = 3;

    public static function apply(int $userId, float $amount): Withdrawal
    {
        return Db::transaction(function () use ($userId, $amount) {
            $rule = PointsRule::order('id', 'desc')->find();
            if (!$rule) {
                throw new \RuntimeException('积分规则未配置');
            }

            if ($amount < (float) $rule->min_withdraw_amount) {
                throw new \RuntimeException('未达到最低提现金额');
            }

            $profile = UserProfile::where('user_id', $userId)->find();
            if (!$profile || !$profile->wechat_qr_url) {
                throw new \RuntimeException('请先绑定微信收款二维码');
            }

            $pointsCost = (int) round($amount * (int) $rule->exchange_rate);

            PointsService::freezeForWithdraw($userId, $pointsCost, '提现冻结');

            return Withdrawal::create([
                'user_id' => $userId,
                'amount' => $amount,
                'points_cost' => $pointsCost,
                'qr_url' => $profile->wechat_qr_url,
                'status' => self::STATUS_PENDING,
            ]);
        });
    }

    public static function approve(int $adminId, int $withdrawalId): void
    {
        Db::transaction(function () use ($adminId, $withdrawalId) {
            $withdrawal = Withdrawal::where('id', $withdrawalId)->lock(true)->findOrFail();
            if ((int) $withdrawal->status !== self::STATUS_PENDING) {
                throw new \RuntimeException('当前状态不可打款');
            }

            PointsService::deductFrozenToWithdrawn((int) $withdrawal->user_id, (int) $withdrawal->points_cost, '提现完成');

            $withdrawal->status = self::STATUS_PAID;
            $withdrawal->reviewed_at = date('Y-m-d H:i:s');
            $withdrawal->admin_id = $adminId;
            $withdrawal->save();

            UserProfile::where('user_id', $withdrawal->user_id)
                ->inc('total_withdrawn', (float) $withdrawal->amount)
                ->update();
        });
    }

    public static function reject(int $adminId, int $withdrawalId, string $reason): void
    {
        Db::transaction(function () use ($adminId, $withdrawalId, $reason) {
            $withdrawal = Withdrawal::where('id', $withdrawalId)->lock(true)->findOrFail();
            if ((int) $withdrawal->status !== self::STATUS_PENDING) {
                throw new \RuntimeException('当前状态不可驳回');
            }

            PointsService::unfreezeToAvailable((int) $withdrawal->user_id, (int) $withdrawal->points_cost, '提现驳回');

            $withdrawal->status = self::STATUS_REJECTED;
            $withdrawal->reviewed_at = date('Y-m-d H:i:s');
            $withdrawal->admin_id = $adminId;
            $withdrawal->reject_reason = $reason;
            $withdrawal->save();
        });
    }
}
