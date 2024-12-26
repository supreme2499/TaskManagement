-- Вставка пользователей
INSERT INTO users (username, password_hash, access_level) VALUES
                                                              ('admin', 'hash123admin', 10),
                                                              ('manager', 'hash123manager', 5),
                                                              ('user1', 'hash123user1', 1),
                                                              ('user2', 'hash123user2', 1),
                                                              ('developer1', 'hash123dev1', 2),
                                                              ('developer2', 'hash123dev2', 2),
                                                              ('qa1', 'hash123qa1', 3),
                                                              ('qa2', 'hash123qa2', 3),
                                                              ('designer1', 'hash123des1', 4),
                                                              ('support1', 'hash123sup1', 4);
-- Вставка задач
INSERT INTO tasks (name_task, description, status, deadline) VALUES
                                                                 ('Fix bug #123', 'Fix the critical bug in the system', 'todo', '2024-12-23 23:59:59'),
                                                                 ('Develop new feature', 'Implement the new feature as discussed', 'in_progress', '2024-12-23 23:59:59'),
                                                                 ('Code review', 'Review the pull requests in the repository', 'done', '2024-12-15 12:00:00'),
                                                                 ('Update documentation', 'Update the project documentation for version 2.0', 'todo', '2025-01-10 18:00:00'),
                                                                 ('Refactor codebase', 'Refactor the old legacy code to improve performance', 'todo', '2024-12-28 09:00:00'),
                                                                 ('Write unit tests', 'Write unit tests for the new module', 'in_progress', '2024-12-23 15:00:00'),
                                                                 ('UI improvements', 'Improve the UI for better user experience', 'todo', '2025-01-05 12:00:00'),
                                                                 ('Server migration', 'Migrate the application to a new server environment', 'todo', '2024-12-23 10:00:00'),
                                                                 ('Security patching', 'Apply critical security patches to the system', 'done', '2024-12-10 11:00:00'),
                                                                 ('Performance testing', 'Test the performance under heavy load', 'in_progress', '2024-12-22 18:00:00');
-- Вставка назначений задач
INSERT INTO task_assignments (user_id, task_id) VALUES
                                                    (1, 1), -- admin assigned to Fix bug #123
                                                    (2, 2), -- manager assigned to Develop new feature
                                                    (3, 3), -- user1 assigned to Code review
                                                    (4, 4), -- user2 assigned to Update documentation
                                                    (5, 5), -- developer1 assigned to Refactor codebase
                                                    (6, 6), -- developer2 assigned to Write unit tests
                                                    (7, 7), -- designer1 assigned to UI improvements
                                                    (8, 8), -- support1 assigned to Server migration
                                                    (9, 9), -- qa1 assigned to Performance testing
                                                    (10, 10), -- support1 assigned to Security patching
                                                    (3, 1), -- user1 also assigned to Fix bug #123
                                                    (4, 2), -- user2 also assigned to Develop new feature
                                                    (5, 3), -- developer1 also assigned to Code review
                                                    (6, 4), -- developer2 also assigned to Update documentation
                                                    (2, 5), -- manager assigned to Refactor codebase
                                                    (3, 6), -- user1 assigned to Write unit tests
                                                    (1, 7), -- admin assigned to UI improvements
                                                    (9, 8), -- qa1 assigned to Server migration
                                                    (7, 9); -- designer1 assigned to Performance testing
