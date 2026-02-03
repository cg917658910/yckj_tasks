<?php

declare(strict_types=1);

namespace app\controller\user;

use app\controller\BaseController;
use app\service\TaskService;
use think\facade\Db;

class ClaimController extends BaseController
{
    public function current()
    {
        $userId = (int) ($this->request->user['user_id'] ?? 0);
        $claim = Db::name('task_claims')->alias('c')
            ->leftJoin('tasks t', 'c.task_id = t.id')
            ->field('c.*, t.title as task_title, t.summary, t.detail, t.doc_url, t.reward_points')
            ->where('c.user_id', $userId)
            ->whereIn('c.status', [TaskService::CLAIM_STATUS_CLAIMED, TaskService::CLAIM_STATUS_SUBMITTED, TaskService::CLAIM_STATUS_REJECTED])
            ->order('c.id', 'desc')
            ->find();

        return json_success(['current' => $claim ?: null]);
    }

    public function history()
    {
        $userId = (int) ($this->request->user['user_id'] ?? 0);
        $list = Db::name('task_claims')->alias('c')
            ->leftJoin('tasks t', 'c.task_id = t.id')
            ->field('c.*, t.title as task_title, t.reward_points')
            ->where('c.user_id', $userId)
            ->order('c.id', 'desc')
            ->select()
            ->toArray();

        return json_success(['list' => $list]);
    }

    public function submit(int $id)
    {
        $userId = (int) ($this->request->user['user_id'] ?? 0);
        $data = [
            'remark' => (string) $this->request->post('remark', ''),
            'images' => $this->request->post('images', []),
        ];
        try {
            validate(\app\validate\TaskSubmitValidate::class)->check($data);
        } catch (\think\exception\ValidateException $e) {
            return json_error($e->getError());
        }

        $remark = $data['remark'];
        $images = $data['images'];

        try {
            TaskService::submitTask($userId, $id, $remark, $images);
        } catch (\Throwable $e) {
            return json_error($e->getMessage());
        }

        return json_success();
    }
}
