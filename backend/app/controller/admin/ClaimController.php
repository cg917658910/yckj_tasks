<?php

declare(strict_types=1);

namespace app\controller\admin;

use app\controller\BaseController;
use app\service\TaskService;
use think\facade\Db;

class ClaimController extends BaseController
{
    public function index()
    {
        $status = $this->request->get('status');
        $query = Db::name('task_claims')->alias('c')
            ->leftJoin('tasks t', 'c.task_id = t.id')
            ->leftJoin('users u', 'c.user_id = u.id')
            ->field('c.*, t.title as task_title, t.reward_points, u.username');

        if ($status !== null && $status !== '') {
            $query->where('c.status', (int) $status);
        }

        $list = $query->order('c.id', 'desc')->select()->toArray();
        return json_success(['list' => $list]);
    }

    public function approve(int $id)
    {
        $rewardPoints = $this->request->post('reward_points');
        $rewardPoints = $rewardPoints === null ? null : (int) $rewardPoints;
        $adminId = (int) ($this->request->admin['admin_id'] ?? 0);

        try {
            TaskService::approveClaim($adminId, $id, $rewardPoints);
        } catch (\Throwable $e) {
            return json_error($e->getMessage());
        }

        return json_success();
    }

    public function reject(int $id)
    {
        $reason = trim((string) $this->request->post('reason', ''));
        if ($reason === '') {
            return json_error('请填写驳回原因');
        }

        $adminId = (int) ($this->request->admin['admin_id'] ?? 0);

        try {
            TaskService::rejectClaim($adminId, $id, $reason);
        } catch (\Throwable $e) {
            return json_error($e->getMessage());
        }

        return json_success();
    }
}
