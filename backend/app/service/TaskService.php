<?php

declare(strict_types=1);

namespace app\service;

use app\model\Task;
use app\model\TaskClaim;
use app\model\TaskSubmission;
use app\model\TaskSubmissionImage;
use think\facade\Db;

class TaskService
{
    public const TASK_STATUS_DRAFT = 0;
    public const TASK_STATUS_ONLINE = 1;
    public const TASK_STATUS_OFFLINE = 2;

    public const CLAIM_STATUS_CLAIMED = 1;
    public const CLAIM_STATUS_SUBMITTED = 2;
    public const CLAIM_STATUS_APPROVED = 3;
    public const CLAIM_STATUS_REJECTED = 4;

    public static function claimTask(int $userId, int $taskId): TaskClaim
    {
        return Db::transaction(function () use ($userId, $taskId) {
            $active = TaskClaim::where('user_id', $userId)
                ->whereIn('status', [self::CLAIM_STATUS_CLAIMED, self::CLAIM_STATUS_SUBMITTED])
                ->lock(true)
                ->find();
            if ($active) {
                throw new \RuntimeException('存在未完成任务，无法领取新任务');
            }

            $task = Task::where('id', $taskId)->lock(true)->find();
            if (!$task || (int) $task->status !== self::TASK_STATUS_ONLINE) {
                throw new \RuntimeException('任务不可领取');
            }

            $exists = TaskClaim::where('task_id', $taskId)->lock(true)->find();
            if ($exists) {
                throw new \RuntimeException('任务已被领取');
            }

            return TaskClaim::create([
                'task_id' => $taskId,
                'user_id' => $userId,
                'status' => self::CLAIM_STATUS_CLAIMED,
                'claimed_at' => date('Y-m-d H:i:s'),
            ]);
        });
    }

    public static function submitTask(int $userId, int $claimId, ?string $remark, array $images): void
    {
        Db::transaction(function () use ($userId, $claimId, $remark, $images) {
            $claim = TaskClaim::where('id', $claimId)->lock(true)->findOrFail();
            if ((int) $claim->user_id !== $userId) {
                throw new \RuntimeException('无权限');
            }
            if (!in_array((int) $claim->status, [self::CLAIM_STATUS_CLAIMED, self::CLAIM_STATUS_REJECTED], true)) {
                throw new \RuntimeException('当前状态不可提交');
            }

            $submission = TaskSubmission::where('claim_id', $claimId)->find();
            if ($submission) {
                $submission->remark = $remark;
                $submission->save();
            } else {
                $submission = TaskSubmission::create([
                    'claim_id' => $claimId,
                    'remark' => $remark,
                ]);
            }

            TaskSubmissionImage::where('submission_id', $submission->id)->delete();
            foreach ($images as $url) {
                TaskSubmissionImage::create([
                    'submission_id' => $submission->id,
                    'image_url' => $url,
                ]);
            }

            $claim->status = self::CLAIM_STATUS_SUBMITTED;
            $claim->submitted_at = date('Y-m-d H:i:s');
            $claim->save();
        });
    }

    public static function approveClaim(int $adminId, int $claimId, ?int $rewardPoints): void
    {
        Db::transaction(function () use ($adminId, $claimId, $rewardPoints) {
            $claim = TaskClaim::where('id', $claimId)->lock(true)->findOrFail();
            if ((int) $claim->status !== self::CLAIM_STATUS_SUBMITTED) {
                throw new \RuntimeException('当前状态不可审核通过');
            }

            $task = Task::where('id', $claim->task_id)->findOrFail();
            $finalPoints = $rewardPoints ?? (int) $task->reward_points;

            $claim->status = self::CLAIM_STATUS_APPROVED;
            $claim->review_result = 1;
            $claim->reviewed_at = date('Y-m-d H:i:s');
            $claim->reward_points_final = $finalPoints;
            $claim->save();

            PointsService::addAvailable((int) $claim->user_id, $finalPoints, 'task_reward', $claim->id, '任务完成奖励');
        });
    }

    public static function rejectClaim(int $adminId, int $claimId, string $reason): void
    {
        Db::transaction(function () use ($adminId, $claimId, $reason) {
            $claim = TaskClaim::where('id', $claimId)->lock(true)->findOrFail();
            if ((int) $claim->status !== self::CLAIM_STATUS_SUBMITTED) {
                throw new \RuntimeException('当前状态不可驳回');
            }

            $claim->status = self::CLAIM_STATUS_REJECTED;
            $claim->review_result = 0;
            $claim->reviewed_at = date('Y-m-d H:i:s');
            $claim->reject_reason = $reason;
            $claim->save();
        });
    }
}
