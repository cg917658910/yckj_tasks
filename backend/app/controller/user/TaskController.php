<?php

declare(strict_types=1);

namespace app\controller\user;

use app\controller\BaseController;
use app\service\TaskService;
use think\facade\Db;

class TaskController extends BaseController
{
    public function index()
    {
        $keyword = trim((string) $this->request->get('keyword', ''));

        $query = Db::name('tasks')->alias('t')
            ->leftJoin('task_claims c', 't.id = c.task_id')
            ->where('t.status', TaskService::TASK_STATUS_ONLINE)
            ->whereNull('c.id')
            ->field('t.*');

        if ($keyword !== '') {
            $query->whereLike('t.title', '%' . $keyword . '%');
        }

        $list = $query->order('t.id', 'desc')->select()->toArray();
        return json_success(['list' => $list]);
    }

    public function detail(int $id)
    {
        $task = Db::name('tasks')->where('id', $id)->find();
        if (!$task) {
            return json_error('任务不存在');
        }
        $claimed = Db::name('task_claims')->where('task_id', $id)->find();
        $task['claimed'] = $claimed ? 1 : 0;
        return json_success($task);
    }

    public function claim(int $id)
    {
        $userId = (int) ($this->request->user['user_id'] ?? 0);
        try {
            $claim = TaskService::claimTask($userId, $id);
        } catch (\Throwable $e) {
            return json_error($e->getMessage());
        }

        return json_success($claim->toArray());
    }
}
